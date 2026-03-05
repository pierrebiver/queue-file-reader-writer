package client

import (
	"bufio"
	"fmt"
	"net"
	"testing"
)

// newTestClient wires a Client to an in-process net.Pipe connection and
// returns the server-side end so tests can simulate server responses.
func newTestClient(t *testing.T) (*Client, net.Conn) {
	t.Helper()
	serverConn, clientConn := net.Pipe()
	t.Cleanup(func() {
		clientConn.Close()
		serverConn.Close()
	})
	c := &Client{
		conn:    clientConn,
		scanner: bufio.NewScanner(clientConn),
	}
	return c, serverConn
}

// fakeServer reads one line from conn then writes response.
func fakeServer(conn net.Conn, response string) {
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	fmt.Fprintf(conn, "%s\n", response)
}

// ── Push ─────────────────────────────────────────────────────────────────────

func TestClient_Push(t *testing.T) {
	c, serverConn := newTestClient(t)
	go fakeServer(serverConn, "OK")

	if err := c.Push("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Push_UnexpectedResponse(t *testing.T) {
	c, serverConn := newTestClient(t)
	go fakeServer(serverConn, "ERR something went wrong")

	if err := c.Push("hello"); err == nil {
		t.Fatal("expected error for unexpected response, got nil")
	}
}

// ── Pop ──────────────────────────────────────────────────────────────────────

func TestClient_Pop(t *testing.T) {
	c, serverConn := newTestClient(t)
	go fakeServer(serverConn, "DATA hello")

	val, err := c.Pop()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "hello" {
		t.Fatalf("expected %q, got %q", "hello", val)
	}
}

func TestClient_Pop_UnexpectedResponse(t *testing.T) {
	c, serverConn := newTestClient(t)
	go fakeServer(serverConn, "UNEXPECTED")

	if _, err := c.Pop(); err == nil {
		t.Fatal("expected error for unexpected response, got nil")
	}
}

// ── SendEOF ──────────────────────────────────────────────────────────────────

func TestClient_SendEOF(t *testing.T) {
	c, serverConn := newTestClient(t)
	go fakeServer(serverConn, "OK")

	if err := c.SendEOF(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
