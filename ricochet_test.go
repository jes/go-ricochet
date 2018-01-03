package goricochet

import (
	"net"
	"testing"
	"time"
)

func SimpleServer() {
	ln, _ := net.Listen("tcp", "127.0.0.1:11000")
	conn, _ := ln.Accept()
	b := make([]byte, 4)
	n, err := conn.Read(b)
	if n == 4 && err == nil {
		conn.Write([]byte{0x01})
	}
	conn.Close()
}

func TestRicochetOpen(t *testing.T) {
	go SimpleServer()
	// Wait for Server to Initialize
	time.Sleep(time.Second)

	rc, err := Open("127.0.0.1:11000|abcdefghijklmno.onion")
	if err == nil {
		if rc.IsInbound {
			t.Errorf("RicochetConnection declares itself as an Inbound connection after an Outbound attempt...that shouldn't happen")
		}
		return
	}
	t.Errorf("RicochetProtocol: Open Failed: %v", err)
}

func BadServer() {
	ln, _ := net.Listen("tcp", "127.0.0.1:11001")
	conn, _ := ln.Accept()
	b := make([]byte, 4)
	n, err := conn.Read(b)
	if n == 4 && err == nil {
		conn.Write([]byte{0xFF})
	}
	conn.Close()
}

func TestRicochetOpenWithError(t *testing.T) {
	go BadServer()
	// Wait for Server to Initialize
	time.Sleep(time.Second)
	_, err := Open("127.0.0.1:11001|abcdefghijklmno.onion")
	if err == nil {
		t.Errorf("Open should have failed because of bad version negotiation.")
	}
}
