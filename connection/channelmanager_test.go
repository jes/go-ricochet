package connection

import (
	"github.com/jes/go-ricochet/channels"
	"testing"
)

type OverrideChatChannel struct {
	channels.ChatChannel
}

// Singleton - for chat channels there can only be one instance per direction
func (cc *OverrideChatChannel) Singleton() bool {
	return false
}

func TestClientManagerDuplicateMultiple(t *testing.T) {
	ccm := NewClientChannelManager()
	chatChannel := new(OverrideChatChannel)
	_, err := ccm.OpenChannelRequestFromPeer(2, chatChannel)
	if err != nil {
		t.Errorf("Opening ChatChannel should have succeeded, instead: %v", err)
	}
	_, err = ccm.OpenChannelRequestFromPeer(4, chatChannel)
	if err != nil {
		t.Errorf("Opening ChatChannel should have succeeded, instead: %v", err)
	}

	channel := ccm.Channel("im.ricochet.chat", channels.Inbound)
	if channel != nil {
		t.Errorf("Get ChatChannel should have failed because there are 2 of them")
	}
}

func TestClientManagerDuplicateChannel(t *testing.T) {
	ccm := NewClientChannelManager()
	chatChannel := new(channels.ChatChannel)
	_, err := ccm.OpenChannelRequestFromPeer(2, chatChannel)
	if err != nil {
		t.Errorf("Opening ChatChannel should have succeeded, instead: %v", err)
	}
	_, err = ccm.OpenChannelRequestFromPeer(2, chatChannel)
	if err == nil {
		t.Errorf("Opening ChatChannel should have failed")
	}

	_, err = ccm.OpenChannelRequestFromPeer(4, chatChannel)
	if err == nil {
		t.Errorf("Opening ChatChannel should have failed because there should be only 1")
	}
}

func TestClientManagerBadServer(t *testing.T) {
	ccm := NewClientChannelManager()
	// Servers are not allowed to open odd numbered channels
	_, err := ccm.OpenChannelRequestFromPeer(3, nil)
	if err == nil {
		t.Errorf("OpenChannelRequestFromPeer should have failed")
	}
}

func TestServerManagerBadClient(t *testing.T) {
	scm := NewServerChannelManager()
	// Clients are not allowed to open even numbered channels
	_, err := scm.OpenChannelRequestFromPeer(2, nil)
	if err == nil {
		t.Errorf("OpenChannelRequestFromPeer should have failed")
	}
}

func TestLocalDuplicate(t *testing.T) {
	scm := NewServerChannelManager()
	chatChannel := new(channels.ChatChannel)
	channel, err := scm.OpenChannelRequest(chatChannel)
	if err != nil {
		t.Errorf("OpenChannelRequest should not have failed: %v", err)
	}

	_, err = scm.OpenChannelRequest(chatChannel)
	if err == nil {
		t.Errorf("OpenChannelRequest should have failed")
	}

	scm.RemoveChannel(channel.ID)
	_, err = scm.OpenChannelRequest(chatChannel)
	if err != nil {
		t.Errorf("OpenChannelRequest should not have failed: %v", err)
	}
}
