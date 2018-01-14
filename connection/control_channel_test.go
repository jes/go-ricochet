package connection

import (
	"github.com/golang/protobuf/proto"
	"github.com/s-rah/go-ricochet/channels"
	"github.com/s-rah/go-ricochet/utils"
	"github.com/s-rah/go-ricochet/wire/control"
	"testing"
)

type MockHandler struct {
	AutoConnectionHandler
}

func (m *MockHandler) GetSupportedChannelTypes() []string {
	return []string{"im.ricochet.chat"}
}

func TestChannelResultNotOpened(t *testing.T) {
	ccm := NewClientChannelManager()
	ctrlChannel := new(ControlChannel)
	ctrlChannel.Init(ccm)
	chatChannel := new(channels.ChatChannel)
	_, err := ccm.OpenChannelRequestFromPeer(2, chatChannel)

	cr := &Protocol_Data_Control.ChannelResult{
		ChannelIdentifier: proto.Int32(2),
		Opened:            proto.Bool(false),
	}
	opened, err := ctrlChannel.ProcessChannelResult(cr)
	if opened != false || err != nil {
		t.Errorf("ProcessChannelResult should have resulted in n channel being opened, and no error %v %v", opened, err)
	}
}

func TestChannelResultError(t *testing.T) {
	ccm := NewClientChannelManager()
	ctrlChannel := new(ControlChannel)
	ctrlChannel.Init(ccm)
	chatChannel := new(channels.ChatChannel)
	_, err := ccm.OpenChannelRequestFromPeer(2, chatChannel)

	cr := &Protocol_Data_Control.ChannelResult{
		ChannelIdentifier: proto.Int32(3),
		Opened:            proto.Bool(false),
	}
	opened, err := ctrlChannel.ProcessChannelResult(cr)
	if opened != false || err != utils.UnexpectedChannelResultError {
		t.Errorf("ProcessChannelResult should have resulted in n channel being opened, and an error %v %v", opened, err)
	}
}

func TestKeepAliveNoResponse(t *testing.T) {
	ctrlChannel := new(ControlChannel)
	ka := &Protocol_Data_Control.KeepAlive{
		ResponseRequested: proto.Bool(false),
	}
	respond, _ := ctrlChannel.ProcessKeepAlive(ka)
	if respond == true {
		t.Errorf("KeepAlive process should have not needed a response %v %v", ka, respond)
	}
}

func TestKeepAliveRequestResponse(t *testing.T) {
	ka := &Protocol_Data_Control.KeepAlive{
		ResponseRequested: proto.Bool(true),
	}
	ctrlChannel := new(ControlChannel)
	respond, _ := ctrlChannel.ProcessKeepAlive(ka)
	if respond == false {
		t.Errorf("KeepAlive process should have produced a response %v %v", ka, respond)
	}
}

func TestEnableFeatures(t *testing.T) {
	handler := new(MockHandler)
	features := []string{"feature1", "im.ricochet.chat"}
	ef := &Protocol_Data_Control.EnableFeatures{
		Feature: features,
	}
	ctrlChannel := new(ControlChannel)
	raw := ctrlChannel.ProcessEnableFeatures(handler, ef)
	res := new(Protocol_Data_Control.Packet)
	err := proto.Unmarshal(raw, res)
	if err != nil || res.GetFeaturesEnabled() == nil {
		t.Errorf("Decoding FeaturesEnabled Packet failed: %v %v", err, res)
	}

	if len(res.GetFeaturesEnabled().GetFeature()) != 1 {
		t.Errorf("Decoding FeaturesEnabled Errored, unexpected length %v", res.GetFeaturesEnabled().GetFeature())
	}

	if res.GetFeaturesEnabled().GetFeature()[0] != "im.ricochet.chat" {
		t.Errorf("Decoding FeaturesEnabled Errored, unexpected feature enabled %v ", res.GetFeaturesEnabled().GetFeature()[0])
	}

}
