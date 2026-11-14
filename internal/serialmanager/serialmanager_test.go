package serialmanager

import (
	"bytes"
	"errors"
	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/internal/protocol"
	"io"
	"strings"
	"testing"
	"time"
)

// mockPort simulates a serial connection without physical hardware.
type mockPort struct {
	incoming bytes.Buffer
	outgoing bytes.Buffer
	closed   bool
	maxRead  int
}

func (p *mockPort) Read(buf []byte) (int, error) {
	if p.closed {
		return 0, errors.New("port closed")
	}

	if p.incoming.Len() == 0 {
		return 0, io.EOF
	}

	if p.maxRead > 0 && len(buf) > p.maxRead {
		buf = buf[:p.maxRead]
	}

	return p.incoming.Read(buf)
}

func (p *mockPort) Write(data []byte) (int, error) {
	if p.closed {
		return 0, errors.New("port closed")
	}

	return p.outgoing.Write(data)
}

func (p *mockPort) Close() error {
	p.closed = true
	return nil
}

func (p *mockPort) SetReadTimeout(time.Duration) error {
	return nil
}

// TestSend checks that messages are sent with a newline.
func TestSend(t *testing.T) {
	port := &mockPort{maxRead: 4}
	manager := New()
	manager.port = port

	err := manager.Send("PING")
	if err != nil {
		t.Fatal(err)
	}

	if got := port.outgoing.String(); got != "PING\n" {
		t.Fatalf("expected PING\\n, got %q", got)
	}
}

// TestReceiveCompleteMessage checks normal message reception.
func TestReceiveCompleteMessage(t *testing.T) {
	port := &mockPort{}
	port.incoming.WriteString("ACCESS_REQUEST\n")

	manager := New()
	manager.port = port

	got, err := manager.Receive()
	if err != nil {
		t.Fatal(err)
	}

	if got != "ACCESS_REQUEST" {
		t.Fatalf("expected ACCESS_REQUEST, got %q", got)
	}
}

// TestReceiveFragmentedMessage checks messages arriving in pieces.
func TestReceiveFragmentedMessage(t *testing.T) {
	port := &mockPort{}
	port.incoming.WriteString("ACCESS_")
	port.incoming.WriteString("REQUEST\n")

	manager := New()
	manager.port = port

	got, err := manager.Receive()
	if err != nil {
		t.Fatal(err)
	}

	if got != "ACCESS_REQUEST" {
		t.Fatalf("expected ACCESS_REQUEST, got %q", got)
	}
}

// TestReceiveMultipleMessages checks that extra messages remain buffered.
func TestReceiveMultipleMessages(t *testing.T) {
	port := &mockPort{}
	port.incoming.WriteString("PING\nSTATUS\n")

	manager := New()
	manager.port = port

	first, err := manager.Receive()
	if err != nil {
		t.Fatal(err)
	}

	second, err := manager.Receive()
	if err != nil {
		t.Fatal(err)
	}

	if first != "PING" || second != "STATUS" {
		t.Fatalf("unexpected messages: %q and %q", first, second)
	}
}

// TestDisconnect checks that the port closes.
func TestDisconnect(t *testing.T) {
	port := &mockPort{}
	manager := New()
	manager.port = port

	err := manager.Disconnect()
	if err != nil {
		t.Fatal(err)
	}

	if !port.closed {
		t.Fatal("expected port to be closed")
	}

	if manager.Connected() {
		t.Fatal("expected manager to be disconnected")
	}
}

// TestSendWithoutConnection checks error handling.
func TestSendWithoutConnection(t *testing.T) {
	manager := New()

	err := manager.Send("PING")
	if err == nil || !strings.Contains(err.Error(), "not connected") {
		t.Fatal("expected not connected error")
	}
}

func TestIdentifyAirLock(t *testing.T) {
	port := &mockPort{}
	port.incoming.WriteString("AIRLOCK:READY:1\n")

	err := identifyAirLock(port)
	if err != nil {
		t.Fatalf("expected successful handshake: %v", err)
	}

	if got := port.outgoing.String(); got != "PING\n" {
		t.Fatalf("expected PING handshake, got %q", got)
	}
}

func TestIdentifyAirLockRejectsUnknownDevice(t *testing.T) {
	port := &mockPort{}
	port.incoming.WriteString("UNKNOWN_DEVICE\n")

	err := identifyAirLock(port)

	if err == nil {
		t.Fatal("expected handshake to fail for an unknown device")
	}
}

func TestSendAuthorizationRequest(t *testing.T) {
	port := &mockPort{}
	manager := New()
	manager.port = port

	request := protocol.AuthorizationRequest{
		Type:      protocol.TypeAuthorizationRequest,
		RequestID: "req-001",
		Action:    "read",
		Resource:  "report.pdf",
	}

	err := manager.SendAuthorizationRequest(request)
	if err != nil {
		t.Fatal(err)
	}

	expected := "{\"type\":\"authorization_request\",\"request_id\":\"req-001\",\"action\":\"read\",\"resource\":\"report.pdf\"}\n"

	if got := port.outgoing.String(); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}
