# 🔌 Section 04-06 — WebSocket: Bidirectional Real-Time Communication

⏱️ Estimated time: **30 minutes**

Welcome to the WebSocket lesson! In this lesson you will learn how to build a **WebSocket server** with Kruda and the `contrib/ws` package -- enabling bidirectional communication between client and server in real time. Perfect for chat, gaming, collaborative editing, and more.

---

## Learning Objectives

- Understand the concept of **WebSocket** and how it differs from regular HTTP
- Learn how to upgrade an HTTP connection to WebSocket using the `contrib/ws` package
- Build a Broadcast Hub pattern for fan-out messages to multiple clients
- Write read/write loops for bidirectional message handling
- Create echo + broadcast handlers that work together

---

## What You Will Learn

By the end of this lesson you will be able to:

- Use `ws.New()` to create an upgrader for WebSocket connections
- Use `upgrader.Upgrade(c, callback)` to upgrade HTTP to WebSocket
- Read messages with `conn.ReadMessage()` and write with `conn.WriteMessage()`
- Build a Hub pattern for broadcasting messages to all clients
- Handle client connect/disconnect correctly
- Use a write pump goroutine to separate read/write operations
- Run WebSocket and REST endpoints together on the same app

---

## Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25+ |
| Git | Latest |
| Text Editor / IDE | VS Code, GoLand, or your preferred editor |
| Terminal | For running `go run` |
| wscat / websocat / Browser | For testing WebSocket |

> If you haven't completed Section 04-05 (MCP CLI), it's recommended to go back and do it first -- [Section 04-05 -- MCP CLI](../05-mcp-server/)

---

## File Structure

```
04-advanced/06-websocket/
├── README.md          <- You are here
├── starter/           <- Starter code (with TODOs to fill in)
│   ├── go.mod
│   └── main.go
└── complete/          <- Complete working solution
    ├── go.mod
    └── main.go
```

- **[starter/](./starter/)** -- Skeleton code that compiles but has `// TODO:` markers for you to complete
- **[complete/](./complete/)** -- Full working solution you can use to compare your answers

---

## What is WebSocket?

**WebSocket** is a protocol that opens a full-duplex (bidirectional) communication channel between client and server over a single TCP connection:

```
┌──────────┐                          ┌──────────────────┐
│  Client  │ -- HTTP Upgrade -------> │   Kruda Server   │
│ (Browser)│ <-- 101 Switching -----  │                  │
│          │                          │   /ws            │
│          │ <----- messages -------> │   Hub (broker)   │
│          │   (bidirectional)        │                  │
└──────────┘                          └──────────────────┘
```

| Feature | WebSocket | SSE | HTTP REST |
|---|---|---|---|
| Direction | Bidirectional | Server -> Client | Request -> Response |
| Connection | Persistent (keep-alive) | Persistent | Closed after response |
| Auto Reconnect | Must implement yourself | Built-in | N/A |
| Best for | Chat, gaming, collaboration | Notifications, feeds | CRUD APIs |
| Protocol | WebSocket (ws:// / wss://) | HTTP | HTTP |

> WebSocket is ideal for cases where both client and server need to send data to each other in real time -- such as a chat room where everyone can send and receive messages simultaneously.

### WebSocket Handshake

A WebSocket connection starts with a regular HTTP request, then upgrades to WebSocket:

```
Client -> Server:
  GET /ws HTTP/1.1
  Upgrade: websocket
  Connection: Upgrade

Server -> Client:
  HTTP/1.1 101 Switching Protocols
  Upgrade: websocket
  Connection: Upgrade
```

> Kruda handles this handshake for you via the `contrib/ws` package -- `upgrader.Upgrade()` will upgrade the connection and return a `*ws.Conn` ready to use.

---

## Architecture: Broadcast Hub Pattern

```
┌──────────┐                           ┌──────────┐
│ Client 1 │ -- conn.ReadMessage() --> │          │ -- conn.WriteMessage() --> Client 2
│          │ <-- conn.WriteMessage() - │   Hub    │ -- conn.WriteMessage() --> Client 3
└──────────┘                           │ (broker) │ -- conn.WriteMessage() --> Client 4
                                       └──────────┘
                                            |
                                       Register / Unregister
                                       per WebSocket connection
```

- **Hub** -- A broker that manages pub/sub between WebSocket clients
- **Client** -- Each WebSocket connection has a `*ws.Conn` and a send channel for outgoing messages
- **Read pump** -- Reads messages with `conn.ReadMessage()` then echoes + broadcasts
- **Write pump** -- A goroutine that reads from the send channel and writes with `conn.WriteMessage()`

---

## Step-by-Step Guide

### Step 1: Open the starter project

```bash
cd 04-advanced/06-websocket/starter
```

Open the `main.go` file -- you will see the skeleton with `// TODO:` comments indicating where to add code.

Note the required imports:

```go
import (
    "github.com/go-kruda/kruda"
    "github.com/go-kruda/kruda/contrib/ws"
)
```

> `contrib/ws` is a WebSocket package separate from core Kruda -- you need to add it to go.mod

### Step 2: Understand the Hub Pattern

The starter includes a `Hub` struct with method signatures already prepared:

- `Register(client)` -- Add a client to the hub
- `Unregister(client)` -- Remove a client from the hub and close the send channel
- `Broadcast(message, sender)` -- Send a message to all clients except the sender
- `ClientCount()` -- Count the number of connected clients

Note that the `Client` struct uses `*ws.Conn` (not `*kruda.WebSocketConn`):

```go
type Client struct {
    conn *ws.Conn
    send chan []byte
}
```

### Step 3: Implement the Broadcast Logic

Replace the `// TODO:` in `Broadcast()`:

```go
func (h *Hub) Broadcast(message []byte, sender *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()

    for client := range h.clients {
        if client == sender {
            continue
        }
        select {
        case client.send <- message:
        default:
        }
    }
}
```

> Use `select` with `default` for non-blocking send -- if a client's buffer is full, it will be skipped instead of blocking.

### Step 4: Write the WebSocket Handler

This is the heart of the lesson! Create an upgrader and WebSocket endpoint in `main()`:

```go
// Create a WebSocket upgrader
upgrader := ws.New(ws.Config{
    MaxMessageSize: 64 * 1024,
})

// Register the WebSocket endpoint
app.Get("/ws", func(c *kruda.Ctx) error {
    return upgrader.Upgrade(c, func(conn *ws.Conn) {
        // Create a client and register with the hub
        client := &Client{
            conn: conn,
            send: make(chan []byte, 64),
        }
        hub.Register(client)
        defer hub.Unregister(client)

        // Notify other clients that someone joined
        hub.Broadcast(
            []byte(fmt.Sprintf(`{"type":"system","message":"a client joined (total: %d)"}`, hub.ClientCount())),
            nil,
        )

        // Write pump -- separate goroutine for sending messages
        go func() {
            for msg := range client.send {
                if err := conn.WriteMessage(ws.TextMessage, msg); err != nil {
                    return
                }
            }
        }()

        // Read pump -- read messages from the client
        for {
            _, msg, err := conn.ReadMessage()
            if err != nil {
                break
            }

            // Echo back to the sender
            echoMsg := fmt.Sprintf(`{"type":"echo","message":%q}`, string(msg))
            if writeErr := conn.WriteMessage(ws.TextMessage, []byte(echoMsg)); writeErr != nil {
                break
            }

            // Broadcast to other clients
            broadcastMsg := fmt.Sprintf(`{"type":"broadcast","from":"client","message":%q}`, string(msg))
            hub.Broadcast([]byte(broadcastMsg), client)
        }

        // Notify other clients that someone left
        hub.Broadcast(
            []byte(fmt.Sprintf(`{"type":"system","message":"a client left (total: %d)"}`, hub.ClientCount())),
            nil,
        )
    })
})
```

> `upgrader.Upgrade(c, callback)` handles the HTTP -> WebSocket upgrade automatically -- the callback function receives a `*ws.Conn` ready to use.

> `conn.ReadMessage()` returns `(messageType int, data []byte, err error)` and `conn.WriteMessage(data []byte)` returns `error`

### Step 5: Add REST Endpoints

Use `kruda.Get[I, O]` for typed REST handlers:

```go
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
```

### Step 6: Run and Test

```bash
go run main.go
```

Test with `wscat` (install with `npm install -g wscat`):

```bash
# Terminal 1 -- connect the first client
wscat -c ws://localhost:3000/ws

# Terminal 2 -- connect the second client
wscat -c ws://localhost:3000/ws
```

Type a message in terminal 1:

```
> Hello!
```

You will see:
- **Terminal 1** (sender): receives an echo message back
- **Terminal 2** (receiver): receives a broadcast message from client 1

Test the REST endpoints:

```bash
# Check the number of connected clients
curl http://localhost:3000/clients

# Health check
curl http://localhost:3000/health
```

Congratulations! You have built a WebSocket server that supports echo + broadcast!

### Test with Browser (WebSocket API)

Open the browser console and run:

```javascript
const ws = new WebSocket("ws://localhost:3000/ws");

ws.onopen = () => {
  console.log("Connected!");
  ws.send("Hello from browser!");
};

ws.onmessage = (e) => {
  console.log("Received:", JSON.parse(e.data));
};

ws.onclose = () => {
  console.log("Disconnected");
};
```

---

## Compare with complete/

If you get stuck, check the solution in **[complete/](./complete/)** and compare with your code, or use the `diff` command:

```bash
diff starter/main.go complete/main.go
```

---

## Key Concepts Summary

| Concept | Description |
|---|---|
| `contrib/ws` package | WebSocket package separate from core Kruda |
| `ws.New(config)` | Create a WebSocket upgrader with configuration |
| `upgrader.Upgrade(c, callback)` | Upgrade HTTP to WebSocket and invoke the callback with `*ws.Conn` |
| `conn.ReadMessage()` | Read a message from a WebSocket client, returns `(msgType, data, err)` |
| `conn.WriteMessage(data)` | Write a message to a WebSocket client |
| `conn.Close(code, reason)` | Close a WebSocket connection |
| Hub pattern | Pub/sub broker for fan-out messages to multiple clients |
| Write pump | Separate goroutine for writing messages to WebSocket |
| Read pump | Main loop that reads messages then echoes + broadcasts |
| Non-blocking send | Use `select` + `default` to prevent a slow client from blocking others |

---

## Next Lesson

Great job! You have learned how to build a bidirectional real-time WebSocket server with Kruda. In the next lesson we will learn **Testing** -- how to write unit tests for Typed Handlers.

[Section 04-07 -- Testing](../07-testing/)
