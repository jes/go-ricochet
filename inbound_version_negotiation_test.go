package goricochet

import (
	"io"
	"net"
	"testing"
	"github.com/s-rah/go-ricochet/utils"
)

func TestBadProtcolLength(t *testing.T) {

	connect := func() {
		conn, err := net.Dial("tcp", ":4000")
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		conn.Write([]byte{0x49, 0x4D})
	}

	l, err := net.Listen("tcp", ":4000")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	go connect()
	conn, err := l.Accept()
	_, err = NegotiateVersionInbound(conn)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("Invalid Error Received. Expected ErrUnexpectedEOF. Got %v", err)
	}

}

func TestNoSupportedVersions(t *testing.T) {

	connect := func() {
		conn, err := net.Dial("tcp", ":4000")
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		conn.Write([]byte{0x49, 0x4D, 0x00, 0xFF})
	}

	l, err := net.Listen("tcp", ":4000")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	go connect()
	conn, err := l.Accept()
	_, err = NegotiateVersionInbound(conn)
	if err != utils.VersionNegotiationError {
		t.Errorf("Invalid Error Received. Expected VersionNegotiationError. Got %v", err)
	}

}

func TestInvalidVersionList(t *testing.T) {

	connect := func() {
		conn, err := net.Dial("tcp", ":4000")
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		conn.Write([]byte{0x49, 0x4D, 0x01})
	}

	l, err := net.Listen("tcp", ":4000")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	go connect()
	conn, err := l.Accept()
	_, err = NegotiateVersionInbound(conn)
	if err != utils.VersionNegotiationError {
		t.Errorf("Invalid Error Received. Expected VersionNegotiationError. Got %v", err)
	}

}

func TestNoCompatibleVersions(t *testing.T) {

	connect := func() {
		conn, err := net.Dial("tcp", ":4000")
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		conn.Write([]byte{0x49, 0x4D, 0x01, 0xFF})
	}

	l, err := net.Listen("tcp", ":4000")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	go connect()
	conn, err := l.Accept()
	_, err = NegotiateVersionInbound(conn)
	if err != utils.VersionNegotiationFailed {
		t.Errorf("Invalid Error Received. Expected VersionNegotiationFailed. Got %v", err)
	}

}
