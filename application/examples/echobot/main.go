package main

import (
	"github.com/s-rah/go-ricochet/application"
	"github.com/s-rah/go-ricochet/channels"
	"github.com/s-rah/go-ricochet/utils"
	"log"
	"time"
)

type EchoBotInstance struct {
	rai *application.ApplicationInstance
	ra  *application.RicochetApplication
}

func (ebi *EchoBotInstance) Init(rai *application.ApplicationInstance, ra *application.RicochetApplication) {
	ebi.rai = rai
	ebi.ra = ra
}

func (ebi *EchoBotInstance) ChatMessage(messageID uint32, when time.Time, message string) bool {
	log.Printf("message from %v - %v", ebi.rai.RemoteHostname, message)
	go ebi.ra.Broadcast(func(rai *application.ApplicationInstance) {
		ebi.SendChatMessage(rai, ebi.rai.RemoteHostname+" "+message)
	})
	return true
}

func (ebi *EchoBotInstance) ChatMessageAck(messageID uint32, accepted bool) {

}

func (ebi *EchoBotInstance) SendChatMessage(rai *application.ApplicationInstance, message string) {
	rai.Connection.Do(func() error {
		// Technically this errors after the second time but we can ignore it.
		rai.Connection.RequestOpenChannel("im.ricochet.chat",
			&channels.ChatChannel{
				Handler: ebi,
			})

		channel := rai.Connection.Channel("im.ricochet.chat", channels.Outbound)
		if channel != nil {
			chatchannel, ok := channel.Handler.(*channels.ChatChannel)
			if ok {
				chatchannel.SendMessage(message)
			}
		}
		return nil
	})
}

func main() {
	echobot := new(application.RicochetApplication)
	pk, err := utils.LoadPrivateKeyFromFile("./testing/private_key")

	if err != nil {
		log.Fatalf("error reading private key file: %v", err)
	}

	l, err := application.SetupOnion("127.0.0.1:9051", "tcp4", "", pk, 9878)

	if err != nil {
		log.Fatalf("error setting up onion service: %v", err)
	}

	af := application.ApplicationInstanceFactory{}
	af.Init()
	af.AddHandler("im.ricochet.chat", func(rai *application.ApplicationInstance) func() channels.Handler {
		ebi := new(EchoBotInstance)
		ebi.Init(rai, echobot)
		return func() channels.Handler {
			chat := new(channels.ChatChannel)
			chat.Handler = ebi
			return chat
		}
	})

	echobot.Init(pk, af, new(application.AcceptAllContactManager))
	log.Printf("echobot listening on %s", l.Addr().String())
	echobot.Run(l)
}
