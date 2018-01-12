package connection

import (
	"github.com/s-rah/go-ricochet/utils"
	"github.com/s-rah/go-ricochet/wire/control"
)

// ProcessEnableFeatures correctly handles a features enabled packet
func ProcessEnableFeatures(handler Handler, ef *Protocol_Data_Control.EnableFeatures) []byte {
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
