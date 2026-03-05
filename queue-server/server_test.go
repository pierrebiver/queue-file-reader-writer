package main

import "testing"

func TestDispatch_Push(t *testing.T) {
	s := NewServer()
	resp := s.dispatch("PUSH hello")
	if resp != "OK" {
		t.Fatalf("expected %q, got %q", "OK", resp)
	}
}

func TestDispatch_Pop(t *testing.T) {
	s := NewServer()
	s.dispatch("PUSH hello")
	resp := s.dispatch("POP")
	if resp != "DATA hello" {
		t.Fatalf("expected %q, got %q", "DATA hello", resp)
	}
}

func TestDispatch_EOF(t *testing.T) {
	s := NewServer()
	resp := s.dispatch("EOF")
	if resp != "OK" {
		t.Fatalf("expected %q, got %q", "OK", resp)
	}
}

func TestDispatch_UnknownCommand(t *testing.T) {
	s := NewServer()
	resp := s.dispatch("FOOBAR")
	if resp != "ERR unknown command" {
		t.Fatalf("expected %q, got %q", "ERR unknown command", resp)
	}
}
