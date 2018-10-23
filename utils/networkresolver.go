package utils

import (
	"golang.org/x/net/proxy"
	"net"
	"strings"
)

const (
	// CannotResolveLocalTCPAddressError is thrown when a local ricochet connection has the wrong format.
	CannotResolveLocalTCPAddressError = Error("CannotResolveLocalTCPAddressError")
	// CannotDialLocalTCPAddressError is thrown when a connection to a local ricochet address fails.
	CannotDialLocalTCPAddressError = Error("CannotDialLocalTCPAddressError")
	// CannotDialRicochetAddressError is thrown when a connection to a ricochet address fails.
	CannotDialRicochetAddressError = Error("CannotDialRicochetAddressError")
)

// NetworkResolver allows a client to resolve various hostnames to connections
// The supported types are onions address are:
//  * ricochet:jlq67qzo6s4yp3sp
//  * jlq67qzo6s4yp3sp
type NetworkResolver struct {
	SOCKSProxy string
}

// Resolve takes a hostname and returns a net.Conn to the derived endpoint
func (nr *NetworkResolver) Resolve(hostname string) (net.Conn, string, error) {
	resolvedHostname := hostname
	if strings.HasPrefix(hostname, "ricochet:") {
		addrParts := strings.Split(hostname, ":")
		resolvedHostname = addrParts[1]
	}

	socksProxy := nr.SOCKSProxy
	if socksProxy == "" {
		socksProxy = "127.0.0.1:9050"
	}
	torDialer, err := proxy.SOCKS5("tcp", socksProxy, nil, proxy.Direct)
	if err != nil {
		return nil, "", err
	}

	conn, err := torDialer.Dial("tcp", resolvedHostname+".onion:9878")
	if err != nil {
		return nil, "", CannotDialRicochetAddressError
	}

	return conn, resolvedHostname, nil
}
