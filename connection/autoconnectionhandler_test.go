package connection

import (
	"github.com/golang/protobuf/proto"
	"github.com/jes/go-ricochet/channels"
	"github.com/jes/go-ricochet/utils"
	"github.com/jes/go-ricochet/wire/control"
	"testing"
)

// Test sending valid packets
func TestInit(t *testing.T) {
	ach := new(AutoConnectionHandler)
	ach.Init()
	ach.RegisterChannelHandler("im.ricochet.auth.hidden-service", func() channels.Handler {
		return &channels.HiddenServiceAuthChannel{}
	})

	// Construct the Open Authentication Channel Message
	messageBuilder := new(utils.MessageBuilder)
	ocm := messageBuilder.OpenAuthenticationChannel(1, [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	// We have just constructed this so there is little
	// point in doing error checking here in the test
	res := new(Protocol_Data_Control.Packet)
	proto.Unmarshal(ocm[:], res)
	opm := res.GetOpenChannel()
	//ocmessage, _ := proto.Marshal(opm)
	handler, err := ach.OnOpenChannelRequest(opm.GetChannelType())

	if err == nil {
		if handler.Type() != "im.ricochet.auth.hidden-service" {
			t.Errorf("Failed to authentication handler: %v", handler.Type())
		}
	} else {
		t.Errorf("Failed to build handler: %v", err)
	}

	types := ach.GetSupportedChannelTypes()
	if len(types) != 1 {
		t.Errorf("Expected only im.ricochet.auth.hidden-service to be supported instead got: %v", types)
	}
}
