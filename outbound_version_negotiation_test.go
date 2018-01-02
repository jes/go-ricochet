package goricochet

import (
	"github.com/s-rah/go-ricochet/utils"
	"net"
	"testing"
	"time"
)

func TestInvalidResponse(t *testing.T) {
	go func() {
		ln, _ := net.Listen("tcp", "127.0.0.1:12000")
		conn, _ := ln.Accept()
		b := make([]byte, 4)
		n, err := conn.Read(b)
		if n == 4 && err == nil {
			conn.Write([]byte{0xFF})
		}
		conn.Close()
	}()
	time.Sleep(time.Second * 1)
	conn, err := net.Dial("tcp", ":12000")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_, err = NegotiateVersionOutbound(conn, "")
	if err != utils.VersionNegotiationFailed {
		t.Errorf("Expected VersionNegotiationFailed got %v", err)
	}
}
