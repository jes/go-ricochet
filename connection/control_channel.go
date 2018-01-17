package connection

import (
	"errors"
	"github.com/s-rah/go-ricochet/utils"
	"github.com/s-rah/go-ricochet/wire/control"
)


// ControlChannel encapsulates logic for the control channel processing
type ControlChannel struct {
	channelManager *ChannelManager
}

// Init sets up a control channel
func (ctrl *ControlChannel) Init(channelManager *ChannelManager) {
	ctrl.channelManager = channelManager
}

// ProcessChannelResult contains the logic for processing a channelresult message
func (ctrl *ControlChannel) ProcessChannelResult(cr *Protocol_Data_Control.ChannelResult) (bool, error) {
	id := cr.GetChannelIdentifier()

	channel, found := ctrl.channelManager.GetChannel(id)

	if !found {
		return false, utils.UnexpectedChannelResultError
	}

	if cr.GetOpened() {
		//rc.traceLog(fmt.Sprintf("channel of type %v opened on %v", channel.Type, id))
		channel.Handler.OpenOutboundResult(nil, cr)
		return true, nil
	}
	//rc.traceLog(fmt.Sprintf("channel of type %v rejected on %v", channel.Type, id))
	channel.Handler.OpenOutboundResult(errors.New(cr.GetCommonError().String()), cr)
	return false, nil
}

// ProcessKeepAlive contains logic for responding to keep alives
func (ctrl *ControlChannel) ProcessKeepAlive(ka *Protocol_Data_Control.KeepAlive) (bool, []byte) {
	if ka.GetResponseRequested() {
		messageBuilder := new(utils.MessageBuilder)
		return true, messageBuilder.KeepAlive(true)
	}
	return false, nil
}

// ProcessEnableFeatures correctly handles a features enabled packet
func (ctrl *ControlChannel) ProcessEnableFeatures(handler Handler, ef *Protocol_Data_Control.EnableFeatures) []byte {
	featuresToEnable := ef.GetFeature()
	supportChannels := handler.GetSupportedChannelTypes()
	result := []string{}
	for _, v := range featuresToEnable {
		for _, s := range supportChannels {
			if v == s {
				result = append(result, v)
			}
		}
	}
	messageBuilder := new(utils.MessageBuilder)
	raw := messageBuilder.FeaturesEnabled(result)
	return raw
}
