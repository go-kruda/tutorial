package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Request / Response Types
// ============================================================

type SendEventInput struct {
	Event string `json:"event"`
	Data  string `json:"data" validate:"required"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// Event Hub -- Pub/Sub for SSE Clients
// ============================================================

type SSEEvent struct {
	Event string
	Data  string
	ID    string
}

type EventHub struct {
	mu      sync.Mutex
	clients map[chan SSEEvent]struct{}
	eventID int
}

func NewEventHub() *EventHub {
	return &EventHub{
		clients: make(map[chan SSEEvent]struct{}),
	}
}

func (h *EventHub) Subscribe() chan SSEEvent {
	h.mu.Lock()
	defer h.mu.Unlock()
	ch := make(chan SSEEvent, 16)
	h.clients[ch] = struct{}{}
	return ch
}

func (h *EventHub) Unsubscribe(ch chan SSEEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, ch)
	close(ch)
}

func (h *EventHub) Broadcast(event, data string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.eventID++
	sseEvent := SSEEvent{
		Event: event,
		Data:  data,
		ID:    fmt.Sprintf("%d", h.eventID),
	}
	for ch := range h.clients {
		select {
		case ch <- sseEvent:
		default:
		}
	}
}

func (h *EventHub) ClientCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.clients)
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	hub := NewEventHub()

	// SSE requires net/http transport because it needs
	// http.Flusher for streaming. Wing does not support this.
	app := kruda.New(kruda.NetHTTP())

	// ── SSE Stream Endpoint ──────────────────────────────────
	//
	// c.SSE() takes a callback that receives an *SSEStream.
	// Inside the callback you can send events and wait for the
	// client to disconnect via stream.Done().
	app.Get("/events", func(c *kruda.Ctx) error {
		ch := hub.Subscribe()
		defer hub.Unsubscribe(ch)

		return c.SSE(func(stream *kruda.SSEStream) error {
			// Send a welcome event.
			stream.Event("connected", map[string]any{
				"message": "connected",
				"clients": hub.ClientCount(),
			})

			// Loop until the client disconnects.
			for {
				select {
				case <-stream.Done():
					return nil

				case event, ok := <-ch:
					if !ok {
						return nil
					}
					stream.EventWithID(event.ID, event.Event, event.Data)
				}
			}
		})
	})

	// ── Send Event Endpoint ──────────────────────────────────
	kruda.Post[SendEventInput, MessageResponse](app, "/send",
		func(c *kruda.C[SendEventInput]) (*MessageResponse, error) {
			eventName := c.In.Event
			if eventName == "" {
				eventName = "message"
			}
			hub.Broadcast(eventName, c.In.Data)
			return &MessageResponse{
				Message: fmt.Sprintf("event '%s' sent to %d client(s)", eventName, hub.ClientCount()),
			}, nil
		},
	)

	// ── Client Count Endpoint ────────────────────────────────
	kruda.Get[struct{}, MessageResponse](app, "/clients",
		func(c *kruda.C[struct{}]) (*MessageResponse, error) {
			return &MessageResponse{
				Message: fmt.Sprintf("%d client(s) connected", hub.ClientCount()),
			}, nil
		},
	)

	// ── Background Heartbeat ─────────────────────────────────
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			hub.Broadcast("heartbeat", fmt.Sprintf(
				`{"time":"%s","clients":%d}`,
				time.Now().Format(time.RFC3339),
				hub.ClientCount(),
			))
		}
	}()

	log.Println("Server starting on :3000 ...")
	log.Println("  SSE stream:   GET  /events")
	log.Println("  Send event:   POST /send")
	log.Println("  Client count: GET  /clients")
	log.Fatal(app.Listen(":3000"))
}
