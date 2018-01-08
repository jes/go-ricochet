package application

import (
	"crypto/rsa"
	"github.com/s-rah/go-ricochet"
	"github.com/s-rah/go-ricochet/connection"
	"github.com/s-rah/go-ricochet/identity"
	"log"
	"net"
	"sync"
)

// RicochetApplication bundles many useful constructs that are
// likely standard in a ricochet application
type RicochetApplication struct {
	contactManager ContactManagerInterface
	privateKey     *rsa.PrivateKey
	l              net.Listener
	instances      []*ApplicationInstance
	lock           sync.Mutex
	aif            ApplicationInstanceFactory
}

func (ra *RicochetApplication) Init(pk *rsa.PrivateKey, af ApplicationInstanceFactory, cm ContactManagerInterface) {
	ra.privateKey = pk
	ra.aif = af
	ra.contactManager = cm
}

func (ra *RicochetApplication) handleConnection(conn net.Conn) {
	rc, err := goricochet.NegotiateVersionInbound(conn)
	if err != nil {
		log.Printf("There was an error")
		conn.Close()
		return
	}

	ich := connection.HandleInboundConnection(rc)

	err = ich.ProcessAuthAsServer(identity.Initialize("", ra.privateKey), ra.contactManager.LookupContact)
	if err != nil {
		log.Printf("There was an error")
		conn.Close()
		return
	}
	rc.TraceLog(true)
	rai := ra.aif.GetApplicationInstance(rc)
	ra.lock.Lock()
	ra.instances = append(ra.instances, rai)
	ra.lock.Unlock()
	rc.Process(rai)
}

func (ra *RicochetApplication) Broadcast(do func(rai *ApplicationInstance)) {
	ra.lock.Lock()
	for _, rai := range ra.instances {
		do(rai)
	}
	ra.lock.Unlock()
}

func (ra *RicochetApplication) Shutdown() {
	log.Printf("Closing")
	ra.l.Close()
	log.Printf("Closed")
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
			log.Printf("Closing")
			return
		}
	}
}
