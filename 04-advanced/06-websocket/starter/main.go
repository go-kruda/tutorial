package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-kruda/kruda"
	"github.com/go-kruda/kruda/contrib/ws"
)

// ============================================================
// Request / Response Types
// ============================================================

// MessageResponse is a generic message envelope.
type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// Broadcast Hub -- Multi-Client Chat
// ============================================================

// Client represents a single WebSocket connection managed by
// the Hub. Each client has a send channel for outgoing messages.
type Client struct {
	conn *ws.Conn
	send chan []byte
}

// Hub manages WebSocket client connections and broadcasts
// messages to all connected clients.
type Hub struct {
	mu      sync.Mutex
	clients map[*Client]struct{}
}

// NewHub creates a Hub ready to accept WebSocket clients.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]struct{}),
	}
}

// Register adds a new client to the hub.
func (h *Hub) Register(client *Client) {
	// TODO: Lock the mutex and add the client to the clients map.
	//
	// Hint:
	//   h.mu.Lock()
	//   defer h.mu.Unlock()
	//   h.clients[client] = struct{}{}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = struct{}{}
}

// Unregister removes a client from the hub and closes its
// send channel.
func (h *Hub) Unregister(client *Client) {
	// TODO: Lock the mutex, check if the client exists, delete
	// it from the map, and close its send channel.
	//
	// Hint:
	//   h.mu.Lock()
	//   defer h.mu.Unlock()
	//   if _, ok := h.clients[client]; ok {
	//       delete(h.clients, client)
	//       close(client.send)
	//   }
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
}

// Broadcast sends a message to every connected client except
// the sender.
func (h *Hub) Broadcast(message []byte, sender *Client) {
	// TODO: Lock the mutex, iterate over all clients, and send
	// the message to each one (except the sender) using a
	// non-blocking send.
	//
	// Hint: Use select with a default case for non-blocking send
	// so a slow client doesn't block delivery to others:
	//
	//   for client := range h.clients {
	//       if client == sender { continue }
	//       select {
	//       case client.send <- message:
	//       default:
	//       }
	//   }
	h.mu.Lock()
	defer h.mu.Unlock()
	_ = message
	_ = sender
}

// ClientCount returns the number of connected clients.
func (h *Hub) ClientCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.clients)
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	// TODO: Create the Broadcast Hub.
	//
	// Example:
	//   hub := NewHub()
	hub := NewHub()

	// WebSocket requires net/http transport (http.Hijacker).
	app := kruda.New(kruda.NetHTTP())

	// TODO: Create a WebSocket upgrader from contrib/ws.
	//
	// Example:
	//   upgrader := ws.New(ws.Config{
	//       MaxMessageSize: 64 * 1024,
	//   })
	_ = ws.New

	// TODO: Register the WebSocket endpoint using app.Get and
	// upgrader.Upgrade(). Inside the upgrade callback you get
	// a *ws.Conn to work with.
	//
	// Steps:
	//   1. app.Get("/ws", func(c *kruda.Ctx) error { ... })
	//   2. Inside: return upgrader.Upgrade(c, func(conn *ws.Conn) { ... })
	//   3. Create a Client{conn: conn, send: make(chan []byte, 64)}
	//   4. hub.Register(client) and defer hub.Unregister(client)
	//   5. Start a write pump goroutine that reads from
	//      client.send and calls conn.WriteMessage(ws.TextMessage, msg)
	//   6. Read loop: conn.ReadMessage() returns (msgType, data, err)
	//      Echo back with conn.WriteMessage(ws.TextMessage, data) and broadcast via
	//      hub.Broadcast()
	//
	// Hint (write pump):
	//   go func() {
	//       for msg := range client.send {
	//           if err := conn.WriteMessage(ws.TextMessage, msg); err != nil { return }
	//       }
	//   }()
	//
	// Hint (read loop):
	//   for {
	//       _, msg, err := conn.ReadMessage()
	//       if err != nil { break }
	//       conn.WriteMessage(ws.TextMessage, []byte(fmt.Sprintf(`{"type":"echo","message":%q}`, string(msg))))
	//       hub.Broadcast([]byte(fmt.Sprintf(`{"type":"broadcast","message":%q}`, string(msg))), client)
	//   }

	// TODO: Register the client count endpoint using kruda.Get.
	//
	// Example:
	//   kruda.Get[struct{}, MessageResponse](app, "/clients", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
	//       return &MessageResponse{Message: fmt.Sprintf("%d client(s) connected", hub.ClientCount())}, nil
	//   })

	// TODO: Register the health check endpoint using kruda.Get.
	//
	// Example:
	//   kruda.Get[struct{}, MessageResponse](app, "/health", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
	//       return &MessageResponse{Message: "WebSocket server is running"}, nil
	//   })

	// Keep imports and variables used to avoid compile errors.
	_ = app
	_ = hub
	_ = fmt.Sprintf

	// TODO: Start the server on port 3000.
	//
	// Example:
	//   log.Fatal(app.Listen(":3000"))
	log.Println("Server starting on :3000 ...")
}
