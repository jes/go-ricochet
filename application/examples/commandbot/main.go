package main

import (
	"github.com/s-rah/go-ricochet/application"
	"github.com/s-rah/go-ricochet/utils"
	"log"
	"time"
)

func main() {
	commandbot := new(application.RicochetApplication)
	pk, err := utils.LoadPrivateKeyFromFile("./testing/private_key")

	if err != nil {
		log.Fatalf("error reading private key file: %v", err)
	}

	l, err := application.SetupOnion("127.0.0.1:9051", "", pk, 9878)

	if err != nil {
		log.Fatalf("error setting up onion service: %v", err)
	}

	commandbot.Init(pk, new(application.AcceptAllContactManager))
	commandbot.OnChatMessage(func(rai *application.RicochetApplicationInstance, id uint32, timestamp time.Time, message string) {
		if message == "/" {
			rai.SendChatMessage(message)
		}
	})
	log.Printf("commandbot listening on %s", l.Addr().String())
	commandbot.Run(l)
}
