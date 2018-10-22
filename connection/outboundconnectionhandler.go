package connection

import (
	"github.com/jes/go-ricochet/channels"
	"github.com/jes/go-ricochet/identity"
	"github.com/jes/go-ricochet/policies"
	"github.com/jes/go-ricochet/utils"
	"sync"
)

// OutboundConnectionHandler is a convieniance wrapper for handling outbound
// connections
type OutboundConnectionHandler struct {
	connection *Connection
}

// HandleOutboundConnection returns an OutboundConnectionHandler given a connection
func HandleOutboundConnection(c *Connection) *OutboundConnectionHandler {
	och := new(OutboundConnectionHandler)
	och.connection = c
	return och
}

// ProcessAuthAsClient blocks until authentication has succeeded or failed with the
// provided identity, or the connection is closed. A non-nil error is returned in all
// cases other than successful authentication.
//
// ProcessAuthAsClient cannot be called at the same time as any other call to a Porcess
// function. Another Process function must be called after this function successfully
// returns to continue handling connection events.
//
// For successful authentication, the `known` return value indicates whether the peer
// accepts us as a known contact. Unknown contacts will generally need to send a contact
// request before any other activity.
func (och *OutboundConnectionHandler) ProcessAuthAsClient(identity identity.Identity) (bool, error) {

	if !identity.Initialized() {
		return false, utils.PrivateKeyNotSetError
	}

	ach := new(AutoConnectionHandler)
	ach.Init()

	// Make sure that calls to Break in this function cannot race
	var breakOnce sync.Once

	var accepted, isKnownContact bool
	authCallback := func(accept, known bool) {
		accepted = accept
		isKnownContact = known
		// Cause the Process() call below to return.
		// If Break() is called from here, it _must_ use go, because this will
		// execute in the Process goroutine, and Break() will deadlock.
		breakOnce.Do(func() { och.connection.Break() })
	}

	processResult := make(chan error, 1)
	go func() {
		// Break Process() if timed out; no-op if Process returned a conn error
		defer func() { breakOnce.Do(func() { och.connection.Break() }) }()
		policy := policies.UnknownPurposeTimeout
		err := policy.ExecuteAction(func() error {
			return och.connection.Process(ach)
		})
		processResult <- err
	}()

	err := och.connection.Do(func() error {
		_, err := och.connection.RequestOpenChannel("im.ricochet.auth.hidden-service",
			&channels.HiddenServiceAuthChannel{
				Identity:         identity,
				ServerHostname:   och.connection.RemoteHostname,
				ClientAuthResult: authCallback,
			})
		return err
	})
	if err != nil {
		breakOnce.Do(func() { och.connection.Break() })
		return false, err
	}

	if err = <-processResult; err != nil {
		return false, err
	}

	if accepted == true {
		return isKnownContact, nil
	}
	return false, utils.ServerRejectedClientConnectionError
}
