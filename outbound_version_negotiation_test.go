package goricochet

import (
	"github.com/jes/go-ricochet/utils"
	"net"
	"testing"
	"time"
)

func TestOutboundVersionNegotiation(t *testing.T) {
	go func() {
		ln, _ := net.Listen("tcp", "127.0.0.1:12001")
		conn, _ := ln.Accept()
		b := make([]byte, 4)
		n, err := conn.Read(b)
		if n == 4 && err == nil {
			conn.Write([]byte{0x01})
		}
		conn.Close()
	}()
	time.Sleep(time.Second * 1)
	conn, err := net.Dial("tcp", ":12001")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_, err = NegotiateVersionOutbound(conn, "")
	if err != nil {
		t.Errorf("Expected success got %v", err)
	}
}

func TestInvalidServer(t *testing.T) {
	go func() {
		ln, _ := net.Listen("tcp", "127.0.0.1:12002")
		conn, _ := ln.Accept()
		b := make([]byte, 4)
		conn.Read(b)
		conn.Write([]byte{})
		conn.Close()
	}()
	time.Sleep(time.Second * 1)
	conn, err := net.Dial("tcp", ":12002")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_, err = NegotiateVersionOutbound(conn, "")
	if err != utils.VersionNegotiationError {
		t.Errorf("Expected VersionNegotiationError got %v", err)
	}
}

func TestInvalidResponse(t *testing.T) {
	go func() {
		ln, _ := net.Listen("tcp", "127.0.0.1:12003")
		conn, _ := ln.Accept()
		b := make([]byte, 4)
		n, err := conn.Read(b)
		if n == 4 && err == nil {
			conn.Write([]byte{0xFF})
		}
		conn.Close()
	}()
	time.Sleep(time.Second * 1)
	conn, err := net.Dial("tcp", ":12003")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_, err = NegotiateVersionOutbound(conn, "")
	if err != utils.VersionNegotiationFailed {
		t.Errorf("Expected VersionNegotiationFailed got %v", err)
	}
}
