package application

import (
	"github.com/s-rah/go-ricochet/channels"
	"github.com/s-rah/go-ricochet/connection"
)

// ApplicationInstance is a  concrete instance of a ricochet application, encapsulating a connection
type ApplicationInstance struct {
	connection.AutoConnectionHandler
	Connection     *connection.Connection
	RemoteHostname string
}

// ApplicationInstanceFactory
type ApplicationInstanceFactory struct {
	handlerMap map[string]func(*ApplicationInstance) func() channels.Handler
}

// Init sets up an Application Factory
func (af *ApplicationInstanceFactory) Init() {
	af.handlerMap = make(map[string]func(*ApplicationInstance) func() channels.Handler)
}

// AddHandler defines a channel type -> handler construct function
func (af *ApplicationInstanceFactory) AddHandler(ctype string, chandler func(*ApplicationInstance) func() channels.Handler) {
	af.handlerMap[ctype] = chandler
}

// GetApplicationInstance buulds a new application instance using a connection as a base.
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

