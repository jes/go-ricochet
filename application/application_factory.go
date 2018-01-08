package application

import (
	"github.com/s-rah/go-ricochet/channels"
	"github.com/s-rah/go-ricochet/connection"
)

// A concrete instance of a ricochet application, encapsulating a connection
type ApplicationInstance struct {
	connection.AutoConnectionHandler
	Connection     *connection.Connection
	RemoteHostname string
}

// Application instance factory
type ApplicationInstanceFactory struct {
	handlerMap map[string]func(*ApplicationInstance) func() channels.Handler
}

// Init setsup an Application Factory
func (af *ApplicationInstanceFactory) Init() {
	af.handlerMap = make(map[string]func(*ApplicationInstance) func() channels.Handler)
}

// AddHandler
func (af *ApplicationInstanceFactory) AddHandler(ctype string, chandler func(*ApplicationInstance) func() channels.Handler) {
	af.handlerMap[ctype] = chandler
}

// GetApplicationInstance,
func (af *ApplicationInstanceFactory) GetApplicationInstance(rc *connection.Connection) *ApplicationInstance {
	rai := new(ApplicationInstance)
	rai.Init()
	rai.RemoteHostname = rc.RemoteHostname
	rai.Connection = rc
	for t, h := range af.handlerMap {
		rai.RegisterChannelHandler(t, h(rai))
	}
	return rai
}
