package testing

import (
	"fmt"
	"github.com/jes/go-ricochet/application"
	"github.com/jes/go-ricochet/channels"
	"github.com/jes/go-ricochet/utils"
	"log"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

type Message struct {
	From, To string
	Message  string
}

type MessageStack interface {
	Add(from, to, message string)
	Get() []Message
}

type Messages struct {
	messages []Message
	lock     sync.Mutex
}

func (messages *Messages) Init() {
	messages.messages = []Message{}
}

func (messages *Messages) Add(from, to, message string) {
	messages.lock.Lock()
	messages.messages = append(messages.messages, Message{from, to, message})
	messages.lock.Unlock()
}

func (messages *Messages) Get() []Message {
	return messages.messages
}

type ChatEchoBot struct {
	onion    string
	rai      *application.ApplicationInstance
	n        int
	Messages MessageStack
}

// We always want bidirectional chat channels
func (bot *ChatEchoBot) OpenInbound() {
	log.Println("OpenInbound() ChatChannel handler called...")
	outboutChatChannel := bot.rai.Connection.Channel("im.ricochet.chat", channels.Outbound)
	if outboutChatChannel == nil {
		bot.rai.Connection.Do(func() error {
			bot.rai.Connection.RequestOpenChannel("im.ricochet.chat",
				&channels.ChatChannel{
					Handler: bot,
				})
			return nil
		})
	}
}

func (bot *ChatEchoBot) ChatMessage(messageID uint32, when time.Time, message string) bool {
	log.Printf("ChatMessage(from: %v, %v", bot.rai.RemoteHostname, message)
	bot.Messages.Add(bot.rai.RemoteHostname, bot.onion, message)
	SendMessage(bot.rai, strconv.Itoa(bot.n)+" witty response")
	bot.n += 1
	return true
}

func SendMessage(rai *application.ApplicationInstance, message string) {
	log.Printf("SendMessage(to: %v, %v)\n", rai.RemoteHostname, message)
	rai.Connection.Do(func() error {

		log.Printf("Finding Chat Channel")
		channel := rai.Connection.Channel("im.ricochet.chat", channels.Outbound)
		if channel != nil {
			log.Printf("Found Chat Channel")
			chatchannel, ok := channel.Handler.(*channels.ChatChannel)
			if ok {
				chatchannel.SendMessage(message)
			}
		} else {
			log.Printf("Could not find chat channel")
		}
		return nil
	})
}

func (bot *ChatEchoBot) ChatMessageAck(messageID uint32, accepted bool) {

}

func TestApplicationIntegration(t *testing.T) {
	startGoRoutines := runtime.NumGoroutine()
	messageStack := &Messages{}
	messageStack.Init()

	fmt.Println("Initializing application factory...")
	af := application.ApplicationInstanceFactory{}
	af.Init()

	af.AddHandler("im.ricochet.contact.request", func(rai *application.ApplicationInstance) func() channels.Handler {
		return func() channels.Handler {
			contact := new(channels.ContactRequestChannel)
			contact.Handler = new(application.AcceptAllContactHandler)
			return contact
		}
	})

	fmt.Println("Starting alice...")
	alice := new(application.RicochetApplication)
	fmt.Println("Generating alice's pk...")
	apk, _ := utils.GeneratePrivateKey()
	aliceAddr, _ := utils.GetOnionAddress(apk)
	fmt.Println("Seting up alice's onion " + aliceAddr + "...")
	al, err := application.SetupOnion("127.0.0.1:9051", "tcp4", "", apk, 9878)
	if err != nil {
		t.Fatalf("Could not setup Onion for Alice: %v", err)
	}

	fmt.Println("Initializing alice...")
	af.AddHandler("im.ricochet.chat", func(rai *application.ApplicationInstance) func() channels.Handler {
		return func() channels.Handler {
			chat := new(channels.ChatChannel)
			chat.Handler = &ChatEchoBot{rai: rai, n: 0, Messages: messageStack, onion: aliceAddr}
			return chat
		}
	})
	alice.Init("Alice", apk, af, new(application.AcceptAllContactManager))
	fmt.Println("Running alice...")
	go alice.Run(al)

	fmt.Println("Starting bob...")
	bob := new(application.RicochetApplication)
	bpk, err := utils.GeneratePrivateKey()
	if err != nil {
		t.Fatalf("Could not setup Onion for Alice: %v", err)
	}
	bobAddr, _ := utils.GetOnionAddress(bpk)
	fmt.Println("Seting up bob's onion " + bobAddr + "...")
	bl, _ := application.SetupOnion("127.0.0.1:9051", "tcp4", "", bpk, 9878)
	af.AddHandler("im.ricochet.chat", func(rai *application.ApplicationInstance) func() channels.Handler {
		return func() channels.Handler {
			chat := new(channels.ChatChannel)
			chat.Handler = &ChatEchoBot{rai: rai, n: 0, Messages: messageStack, onion: bobAddr}
			return chat
		}
	})
	bob.Init("Bob", bpk, af, new(application.AcceptAllContactManager))
	go bob.Run(bl)

	fmt.Println("Waiting for alice and bob hidden services to percolate...")
	time.Sleep(60 * time.Second)
	runningGoRoutines := runtime.NumGoroutine()

	fmt.Println("Alice connecting to Bob...")
	// out going rc from alice to bob
	alicei, err := alice.Open(bobAddr, "It's alice")
	if err != nil {
		t.Fatalf("Error Alice connecting to Bob: %v", err)
	}
	time.Sleep(10 * time.Second)

	fmt.Println("Alice request open chat channel...")
	// TODO: opening a channel should be easier?
	alicei.Connection.Do(func() error {
		handler, err := alicei.OnOpenChannelRequest("im.ricochet.chat")
		if err != nil {
			log.Printf("Could not get chat handler!\n")
			return err
		}
		_, err = alicei.Connection.RequestOpenChannel("im.ricochet.chat", handler)
		return err
	})
	time.Sleep(5 * time.Second)

	fmt.Println("Alice sending message to Bob...")
	SendMessage(alicei, "Hello Bob!")

	if err != nil {
		log.Fatal("Error dialing from Alice to Bob: ", err)
	}

	time.Sleep(10 * time.Second)

	// should now be connected to bob
	connectedGoRoutines := runtime.NumGoroutine()

	fmt.Println("Shutting bob down...")
	bob.Shutdown()

	time.Sleep(15 * time.Second)

	bobShutdownGoRoutines := runtime.NumGoroutine()

	fmt.Println("Shutting alice down...")
	alice.Shutdown()
	time.Sleep(15 * time.Second)

	finalGoRoutines := runtime.NumGoroutine()

	fmt.Printf("startGoRoutines: %v\nrunningGoROutines: %v\nconnectedGoRoutines: %v\nBobShutdownGoRoutines: %v\nfinalGoRoutines: %v\n", startGoRoutines, runningGoRoutines, connectedGoRoutines, bobShutdownGoRoutines, finalGoRoutines)

	if bobShutdownGoRoutines != startGoRoutines+1 {
		t.Errorf("After shutting down bob, go routines were not start + 1 (alice) value. Expected: %v Actual %v", startGoRoutines+1, bobShutdownGoRoutines)
	}

	if finalGoRoutines != startGoRoutines {
		t.Errorf("After shutting alice and bob down, go routines were not at start value. Expected: %v Actual: %v", startGoRoutines, finalGoRoutines)
	}

	fmt.Println("Messages:")
	for _, message := range messageStack.Get() {
		fmt.Printf("  from:%v to:%v '%v'\n", message.From, message.To, message.Message)
	}

	messages := messageStack.Get()
	if messages[0].Message != "Hello Bob!" || messages[1].Message != "0 witty response" {
		t.Errorf("Message history did not contain first two expected messages!")
	}
}
