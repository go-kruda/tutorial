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
//
// These structs define the JSON payloads for our REST endpoints.
// The WebSocket messages use their own format -- plain text or
// JSON strings sent over the WebSocket connection.

// MessageResponse is a generic message envelope.
type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// Broadcast Hub -- Multi-Client Chat
// ============================================================
//
// The Hub manages all active WebSocket connections and broadcasts
// messages to every connected client. This is the classic
// "chat room" pattern:
//
//   1. A client connects via WebSocket -> Hub registers it
//   2. A client sends a message -> Hub broadcasts to all others
//   3. A client disconnects -> Hub removes it
//
// Why a hub pattern?
// ------------------
// WebSocket is bidirectional -- both client and server can send
// messages at any time. When building multi-client features
// (chat, collaborative editing, live dashboards), you need a
// central broker to fan-out messages. The hub pattern keeps
// this logic clean and separate from your WebSocket handlers.

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
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = struct{}{}
}

// Unregister removes a client from the hub and closes its
// send channel. Always call this when a client disconnects
// to avoid goroutine and channel leaks.
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
}

// Broadcast sends a message to every connected client except
// the sender. This is the classic chat broadcast -- when one
// client sends a message, all others receive it.
//
// The sender parameter can be nil to broadcast to everyone
// (useful for system messages like join/leave notifications).
func (h *Hub) Broadcast(message []byte, sender *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		// Skip the sender so they don't receive their own message.
		if client == sender {
			continue
		}

		// Non-blocking send: if the client's buffer is full,
		// we skip rather than block the broadcaster. In a
		// production app you might want to disconnect slow
		// clients instead.
		select {
		case client.send <- message:
		default:
		}
	}
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
	// ── 1. Create the Broadcast Hub ──────────────────────────
	//
	// The hub is the central broker for all WebSocket clients.
	// When a client sends a message, the hub fans it out to
	// every other connected client.
	hub := NewHub()

	// ── 2. Create the Kruda Application ──────────────────────
	//
	// WebSocket upgrade requires net/http transport because it
	// needs http.Hijacker. Wing and fasthttp do not support this.
	app := kruda.New(kruda.NetHTTP())

	// ── 3. Create the WebSocket Upgrader ─────────────────────
	//
	// The upgrader from contrib/ws handles the HTTP -> WebSocket
	// upgrade (101 Switching Protocols handshake). You configure
	// it once and reuse it for all WebSocket endpoints.
	upgrader := ws.New(ws.Config{
		MaxMessageSize: 64 * 1024,
	})

	// ── 4. Register the WebSocket Endpoint ───────────────────
	//
	// GET /ws is the WebSocket endpoint. Clients connect here
	// using the browser's WebSocket API or a CLI tool like
	// websocat / wscat:
	//
	//   const ws = new WebSocket("ws://localhost:3000/ws");
	//   ws.onmessage = (e) => console.log(e.data);
	//   ws.send("Hello!");
	//
	// upgrader.Upgrade() handles the HTTP -> WebSocket upgrade
	// automatically and gives you a *ws.Conn to work with.
	app.Get("/ws", func(c *kruda.Ctx) error {
		return upgrader.Upgrade(c, func(conn *ws.Conn) {
			// Create a client and register it with the hub.
			client := &Client{
				conn: conn,
				send: make(chan []byte, 64),
			}
			hub.Register(client)
			defer hub.Unregister(client)

			// Notify all other clients that someone joined.
			hub.Broadcast(
				[]byte(fmt.Sprintf(`{"type":"system","message":"a client joined (total: %d)"}`, hub.ClientCount())),
				nil,
			)

			// ── Write pump ───────────────────────────────────
			//
			// A separate goroutine reads from the client's send
			// channel and writes messages to the WebSocket. This
			// decouples the broadcast fan-out from the read loop.
			go func() {
				for msg := range client.send {
					if err := conn.WriteMessage(ws.TextMessage, msg); err != nil {
						return
					}
				}
			}()

			// ── Read pump ────────────────────────────────────
			//
			// The main goroutine reads messages from the WebSocket.
			// For each message:
			//   1. Echo it back to the sender
			//   2. Broadcast the original message to all others
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					// Client disconnected or error -- exit the loop.
					break
				}

				// Echo the message back to the sender.
				echoMsg := fmt.Sprintf(`{"type":"echo","message":%q}`, string(msg))
				if writeErr := conn.WriteMessage(ws.TextMessage, []byte(echoMsg)); writeErr != nil {
					break
				}

				// Broadcast the original message to all other clients.
				broadcastMsg := fmt.Sprintf(`{"type":"broadcast","from":"client","message":%q}`, string(msg))
				hub.Broadcast([]byte(broadcastMsg), client)
			}

			// Notify others that this client left.
			hub.Broadcast(
				[]byte(fmt.Sprintf(`{"type":"system","message":"a client left (total: %d)"}`, hub.ClientCount())),
				nil,
			)
		})
	})

	// ── 5. Register REST Endpoints ───────────────────────────
	//
	// WebSocket and REST routes coexist on the same app. These
	// endpoints are regular HTTP -- they don't use WebSocket.
	kruda.Get[struct{}, MessageResponse](app, "/clients", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
		return &MessageResponse{
			Message: fmt.Sprintf("%d client(s) connected", hub.ClientCount()),
		}, nil
	})

	kruda.Get[struct{}, MessageResponse](app, "/health", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
		return &MessageResponse{
			Message: "WebSocket server is running",
		}, nil
	})

	// ── 6. Start the Server ──────────────────────────────────
	log.Println("WebSocket server starting on :3000 ...")
	log.Println("Endpoints:")
	log.Println("  GET  /ws       - WebSocket endpoint")
	log.Println("  GET  /clients  - Connected client count")
	log.Println("  GET  /health   - Health check")
	log.Fatal(app.Listen(":3000"))
}
