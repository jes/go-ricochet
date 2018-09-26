package connection

import (
	"github.com/jes/go-ricochet/channels"
	"github.com/jes/go-ricochet/utils"
)

// AutoConnectionHandler implements the ConnectionHandler interface on behalf of
// the provided application type by automatically providing support for any
// built-in channel type whose high level interface is implemented by the
// application. For example, if the application's type implements the
// ChatChannelHandler interface, `im.ricochet.chat` will be available to the peer.
//
// The application handler can be any other type. To override or augment any of
// AutoConnectionHandler's behavior (such as adding new channel types, or reacting
// to connection close events), this type can be embedded in the type that it serves.
type AutoConnectionHandler struct {
	handlerMap      map[string]func() channels.Handler
	connection      *Connection
	OnClosedHandler func(error)
}

// Init ...
// TODO: Split this into client and server init
func (ach *AutoConnectionHandler) Init() {
	ach.handlerMap = make(map[string]func() channels.Handler)
}

// OnReady ...
func (ach *AutoConnectionHandler) OnReady(oc *Connection) {
	ach.connection = oc
}

// OnClosed is called when the OpenConnection has closed for any reason.
func (ach *AutoConnectionHandler) OnClosed(err error) {
	if ach.OnClosedHandler != nil {
		ach.OnClosedHandler(err)
	}
}

// GetSupportedChannelTypes returns a list of channel types that are registered with the handler.
func (ach *AutoConnectionHandler) GetSupportedChannelTypes() []string {
	supported := []string{}
	for k := range ach.handlerMap {
		supported = append(supported, k)
	}
	return supported
}

// RegisterChannelHandler ...
func (ach *AutoConnectionHandler) RegisterChannelHandler(ctype string, handler func() channels.Handler) {
	_, exists := ach.handlerMap[ctype]
	if !exists {
		ach.handlerMap[ctype] = handler
	}
}

// OnOpenChannelRequest ...
func (ach *AutoConnectionHandler) OnOpenChannelRequest(ctype string) (channels.Handler, error) {
	handler, ok := ach.handlerMap[ctype]
	if ok {
		h := handler()
		return h, nil
	}
	return nil, utils.UnknownChannelTypeError
}
