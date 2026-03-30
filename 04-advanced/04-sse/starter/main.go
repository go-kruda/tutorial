package main

import (
	"fmt"
	"log"
	"sync"

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
// Event Hub -- Pub/Sub (provided)
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
	return &EventHub{clients: make(map[chan SSEEvent]struct{})}
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
	sseEvent := SSEEvent{Event: event, Data: data, ID: fmt.Sprintf("%d", h.eventID)}
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

	// SSE requires net/http transport (http.Flusher).
	app := kruda.New(kruda.NetHTTP())

	// TODO: implement SSE stream endpoint GET /events
	//
	// Hint: Use c.SSE() with the callback pattern
	//
	//   app.Get("/events", func(c *kruda.Ctx) error {
	//       ch := hub.Subscribe()
	//       defer hub.Unsubscribe(ch)
	//
	//       return c.SSE(func(stream *kruda.SSEStream) error {
	//           stream.Event("connected", map[string]any{"message": "connected"})
	//
	//           for {
	//               select {
	//               case <-stream.Done():
	//                   return nil
	//               case event, ok := <-ch:
	//                   if !ok { return nil }
	//                   stream.EventWithID(event.ID, event.Event, event.Data)
	//               }
	//           }
	//       })
	//   })

	// TODO: implement POST /send -- send an event to all SSE clients
	//
	// Hint: Use kruda.Post[SendEventInput, MessageResponse]
	//   hub.Broadcast(eventName, c.In.Data)

	// TODO: implement GET /clients -- show the number of connected clients
	//
	// Hint: Use kruda.Get[struct{}, MessageResponse]

	_ = hub
	log.Println("Server starting on :3000 ...")
	log.Fatal(app.Listen(":3000"))
}
