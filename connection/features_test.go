package connection

import (
	"github.com/golang/protobuf/proto"
	"github.com/s-rah/go-ricochet/wire/control"
	"testing"
)

type MockHandler struct {
	AutoConnectionHandler
}

func (m *MockHandler) GetSupportedChannelTypes() []string {
	return []string{"im.ricochet.chat"}
}

func TestEnableFeatures(t *testing.T) {
	handler := new(MockHandler)
	features := []string{"feature1", "im.ricochet.chat"}
	ef := &Protocol_Data_Control.EnableFeatures{
		Feature: features,
	}
	raw := ProcessEnableFeatures(handler, ef)
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
