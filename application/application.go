package application

import (
	"crypto/rsa"
	"github.com/jes/go-ricochet"
	"github.com/jes/go-ricochet/channels"
	"github.com/jes/go-ricochet/connection"
	"github.com/jes/go-ricochet/identity"
	"log"
	"net"
	"sync"
)

// RicochetApplication bundles many useful constructs that are
// likely standard in a ricochet application
type RicochetApplication struct {
	contactManager     ContactManagerInterface
	privateKey         *rsa.PrivateKey
	name               string
	l                  net.Listener
	instances          []*ApplicationInstance
	lock               sync.Mutex
	aif                ApplicationInstanceFactory
	OnNewPeer          func(*ApplicationInstance, string)
	OnAuthenticated    func(*ApplicationInstance, bool)
	MakeContactHandler func(*ApplicationInstance) channels.ContactRequestChannelHandler
	SOCKSProxy         string
}

func (ra *RicochetApplication) Init(name string, pk *rsa.PrivateKey, af ApplicationInstanceFactory, cm ContactManagerInterface) {
	ra.name = name
	ra.privateKey = pk
	ra.aif = af
	ra.contactManager = cm
	ra.SOCKSProxy = "127.0.0.1:9050"
	ra.MakeContactHandler = func(*ApplicationInstance) channels.ContactRequestChannelHandler {
		return new(AcceptAllContactHandler)
	}
}

// TODO: Reimplement OnJoin, OnLeave Events.
func (ra *RicochetApplication) handleConnection(conn net.Conn) {
	rc, err := goricochet.NegotiateVersionInbound(conn)
	if err != nil {
		log.Printf("There was an error")
		conn.Close()
		return
	}

	ich := connection.HandleInboundConnection(rc)

	rai := ra.aif.GetApplicationInstance(rc)
	lookupContactFn := func(hostname string, publicKey rsa.PublicKey) (bool, bool) {
		rai.RemoteHostname = hostname // XXX: without this here, I think RemoteHostname only gets set for outgoing connections
		if ra.OnNewPeer != nil {
			ra.OnNewPeer(rai, hostname)
		}
		return ra.contactManager.LookupContact(hostname, publicKey)
	}
	err = ich.ProcessAuthAsServer(identity.Initialize(ra.name, ra.privateKey), lookupContactFn)
	if err != nil {
		log.Printf("There was an error")
		conn.Close()
		return
	}
	rc.TraceLog(true)
	ra.lock.Lock()
	ra.instances = append(ra.instances, rai)
	ra.lock.Unlock()
	rc.Process(rai)
}

func (ra *RicochetApplication) HandleApplicationInstance(rai *ApplicationInstance) {
	ra.lock.Lock()
	ra.instances = append(ra.instances, rai)
	ra.lock.Unlock()
}

// Open a connection to another Ricochet peer at onionAddress. If they are unknown to use, use requestMessage (otherwise can be blank)
func (ra *RicochetApplication) Open(onionAddress string, requestMessage string) (*ApplicationInstance, error) {
	rc, err := goricochet.OpenWithProxy(onionAddress, ra.SOCKSProxy)
	if err != nil {
		log.Printf("Error in application.Open(): %v\n", err)
		return nil, err
	}
	rc.TraceLog(true)

	known, err := connection.HandleOutboundConnection(rc).ProcessAuthAsClient(identity.Initialize(ra.name, ra.privateKey))
	rai := ra.aif.GetApplicationInstance(rc)
	if ra.OnAuthenticated != nil {
		ra.OnAuthenticated(rai, known)
	}
	go rc.Process(rai)

	if !known {
		err := rc.Do(func() error {
			_, err := rc.RequestOpenChannel("im.ricochet.contact.request",
				&channels.ContactRequestChannel{
					Handler: ra.MakeContactHandler(rai),
					Name:    ra.name,
					Message: requestMessage,
				})
			return err
		})
		if err != nil {
			log.Printf("could not contact %s", err)
		}
	}

	ra.HandleApplicationInstance(rai)
	return rai, nil
}

func (ra *RicochetApplication) Broadcast(do func(rai *ApplicationInstance)) {
	ra.lock.Lock()
	for _, rai := range ra.instances {
		do(rai)
	}
	ra.lock.Unlock()
}

func (ra *RicochetApplication) Shutdown() {
	ra.l.Close()
	for _, instance := range ra.instances {
		instance.Connection.Conn.Close()
	}
}

func (ra *RicochetApplication) Run(l net.Listener) {
	if ra.privateKey == nil || ra.contactManager == nil {
		return
	}
	ra.l = l
	var err error
	for err == nil {
		conn, err := ra.l.Accept()
		if err == nil {
			go ra.handleConnection(conn)
		} else {
			return
		}
	}
}
