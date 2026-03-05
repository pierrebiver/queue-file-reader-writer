package command

import "testing"

// mockStore uses a buffered channel to satisfy the blocking Pop() contract.
type mockStore struct {
	ch chan string
}

func newMockStore(cap int) *mockStore {
	return &mockStore{ch: make(chan string, cap)}
}

func (m *mockStore) Push(item string) {
	m.ch <- item
}

func (m *mockStore) Pop() string {
	return <-m.ch
}

// ── PushCommand ──────────────────────────────────────────────────────────────

func TestPushCommand_Is_MatchingLine(t *testing.T) {
	cmd := PushCommand{}
	if !cmd.Is("PUSH hello") {
		t.Fatal("expected Is to return true for 'PUSH hello'")
	}
}

func TestPushCommand_Is_NonMatchingLine(t *testing.T) {
	cmd := PushCommand{}
	if cmd.Is("POP") {
		t.Fatal("expected Is to return false for 'POP'")
	}
}

func TestPushCommand_Is_MissingSpace(t *testing.T) {
	cmd := PushCommand{}
	if cmd.Is("PUSH") {
		t.Fatal("expected Is to return false for 'PUSH' without trailing space")
	}
}

func TestPushCommand_CommandName(t *testing.T) {
	cmd := PushCommand{}
	if got := cmd.CommandName(); got != "PUSH" {
		t.Fatalf("expected %q, got %q", "PUSH", got)
	}
}

func TestPushCommand_Execute_ValidPayload(t *testing.T) {
	store := newMockStore(1)
	cmd := PushCommand{}
	resp, err := cmd.Execute("PUSH hello", store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "OK" {
		t.Fatalf("expected %q, got %q", "OK", resp)
	}
	if got := store.Pop(); got != "hello" {
		t.Fatalf("expected %q in store, got %q", "hello", got)
	}
}

// ── PopCommand ───────────────────────────────────────────────────────────────

func TestPopCommand_Is_MatchingLine(t *testing.T) {
	cmd := PopCommand{}
	if !cmd.Is("POP") {
		t.Fatal("expected Is to return true for 'POP'")
	}
}

func TestPopCommand_Is_NonMatchingLine(t *testing.T) {
	cmd := PopCommand{}
	if cmd.Is("PUSH foo") {
		t.Fatal("expected Is to return false for 'PUSH foo'")
	}
}

func TestPopCommand_CommandName(t *testing.T) {
	cmd := PopCommand{}
	if got := cmd.CommandName(); got != "POP" {
		t.Fatalf("expected %q, got %q", "POP", got)
	}
}

func TestPopCommand_Execute_QueueHasItem(t *testing.T) {
	store := newMockStore(1)
	store.Push("hello")
	cmd := PopCommand{}
	resp, err := cmd.Execute("POP", store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "DATA hello" {
		t.Fatalf("expected %q, got %q", "DATA hello", resp)
	}
}

// ── EOFCommand ───────────────────────────────────────────────────────────────

func TestEOFCommand_Is_MatchingLine(t *testing.T) {
	cmd := EOFCommand{}
	if !cmd.Is("EOF") {
		t.Fatal("expected Is to return true for 'EOF'")
	}
}

func TestEOFCommand_Is_NonMatchingLine(t *testing.T) {
	cmd := EOFCommand{}
	if cmd.Is("POP") {
		t.Fatal("expected Is to return false for 'POP'")
	}
}

func TestEOFCommand_CommandName(t *testing.T) {
	cmd := EOFCommand{}
	if got := cmd.CommandName(); got != "EOF" {
		t.Fatalf("expected %q, got %q", "EOF", got)
	}
}

func TestEOFCommand_Execute_ReturnsOK(t *testing.T) {
	store := newMockStore(1)
	cmd := EOFCommand{}
	resp, err := cmd.Execute("EOF", store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "OK" {
		t.Fatalf("expected %q, got %q", "OK", resp)
	}
}

func TestEOFCommand_Execute_PushesSentinel(t *testing.T) {
	store := newMockStore(1)
	cmd := EOFCommand{}
	cmd.Execute("EOF", store)
	if got := store.Pop(); got != "EOF" {
		t.Fatalf("expected sentinel %q in store, got %q", "EOF", got)
	}
}

// ── Registry ─────────────────────────────────────────────────────────────────

func TestRegistry_ContainsAllCommands(t *testing.T) {
	names := map[string]bool{}
	for _, cmd := range Registry {
		names[cmd.CommandName()] = true
	}
	for _, want := range []string{"PUSH", "POP", "EOF"} {
		if !names[want] {
			t.Errorf("Registry missing command %q", want)
		}
	}
}

func TestRegistry_NoAmbiguity(t *testing.T) {
	probes := []string{"PUSH hello", "POP", "EOF"}
	for _, line := range probes {
		matches := 0
		for _, cmd := range Registry {
			if cmd.Is(line) {
				matches++
			}
		}
		if matches != 1 {
			t.Errorf("line %q matched %d commands, want exactly 1", line, matches)
		}
	}
}
