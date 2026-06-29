package main

import (
	"strings"
	"testing"

	"github.com/go-kruda/kruda"
)

// TestSnapshotSSE exercises the finite /snapshot SSE endpoint with the
// in-memory TestClient.SSE helper. TestClient.SSE runs the handler to
// completion, so it works on a finite stream like /snapshot — it would
// hang on the infinite /events stream (which loops until disconnect).
func TestSnapshotSSE(t *testing.T) {
	hub := NewEventHub()
	app := newApp(hub)
	client := kruda.NewTestClient(app)

	result, err := client.SSE("/snapshot")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != 200 {
		t.Fatalf("expected status 200, got %d", result.Status)
	}
	if len(result.Events) != 1 {
		t.Fatalf("expected 1 event, got %d: %v", len(result.Events), result.Events)
	}
	if !strings.Contains(result.Events[0], "event: snapshot") {
		t.Errorf("expected a 'snapshot' event, got: %q", result.Events[0])
	}
	if !strings.Contains(result.Events[0], "clients") {
		t.Errorf("expected client count in event data, got: %q", result.Events[0])
	}
}
