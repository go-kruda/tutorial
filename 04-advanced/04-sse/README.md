# 📡 Section 04-04 — SSE: Real-Time Server-Sent Events

⏱️ Estimated time: **30 minutes**

Welcome to the Server-Sent Events (SSE) lesson! 🎉 In this lesson you'll learn how to **push real-time data** from server to client using Kruda's SSE — perfect for notifications, live feeds, dashboards, and anything else that needs instant updates.

---

## 🎯 Lesson Objectives

- Understand the concept of **Server-Sent Events (SSE)** and how it differs from WebSocket
- Learn how to open an SSE stream with `c.SSE(func(stream *kruda.SSEStream) error {...})`
- Build an Event Hub (pub/sub) for broadcasting events to multiple clients
- Send events with event name, data, and ID
- Use heartbeat to keep the connection active

---

## 📚 What You'll Learn

By the end of this lesson you'll be able to:

- ✅ Open an SSE stream with `c.SSE(func(stream *kruda.SSEStream) error {...})` callback pattern
- ✅ Send events with `stream.Event(name, data)` and `stream.EventWithID(id, name, data)`
- ✅ Build an Event Hub pattern for fan-out events to multiple clients
- ✅ Handle client connect/disconnect properly
- ✅ Use heartbeat to prevent connection timeout
- ✅ Trigger SSE events from a regular REST endpoint
- ✅ Use the browser EventSource API to receive events

---

## 📋 Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25 or higher |
| Git | Latest version |
| Text Editor / IDE | VS Code, GoLand, or your preferred editor |
| Terminal | For running `go run` |
| Browser / curl | For testing the SSE stream |

> 💡 If you haven't completed Section 04-03 (OpenAPI) yet, it's recommended to go back and do it first — 👉 [Section 04-03 — OpenAPI Generator](../03-openapi/)

---

## 📂 File Structure

```
04-advanced/04-sse/
├── README.md          ← 📖 You are here
├── starter/           ← 🏗️ Starter code (with TODOs to fill in)
│   ├── go.mod
│   └── main.go
└── complete/          ← ✅ Fully working reference implementation
    ├── go.mod
    └── main.go
```

- 📁 **[starter/](./starter/)** — A compilable skeleton with `// TODO:` placeholders for you to fill in
- 📁 **[complete/](./complete/)** — A fully working reference implementation you can compare against

---

## 💡 What is SSE?

**Server-Sent Events (SSE)** is a protocol that lets the server push data to the client over a single HTTP connection:

| Feature | SSE | WebSocket |
|---|---|---|
| Direction | Server → Client (one-way) | Bidirectional |
| Protocol | Standard HTTP | WebSocket protocol |
| Auto reconnect | ✅ Built-in | ❌ Must implement yourself |
| Event types | ✅ Supports named events | ❌ No built-in support |
| Ease of use | ✅ Very easy | ⚠️ More complex |
| Best for | Notifications, feeds, dashboards | Chat, gaming, real-time collaboration |

> 📡 SSE is ideal for cases where the server needs to push data to the client but the client doesn't need to send data back — much simpler than WebSocket!

### SSE Event Format

SSE uses a simple plain text format:

```
event: notification
data: {"message":"Hello!"}
id: 1

event: heartbeat
data: {"time":"2025-01-01T00:00:00Z"}
id: 2

```

Each event consists of:
- `event:` — Event name (optional, defaults to "message")
- `data:` — Event payload
- `id:` — ID for reconnection (optional)

> ✨ Kruda handles this format for you automatically via `stream.Event()` / `stream.EventWithID()` — no need to write the format yourself

---

## 🏗️ Architecture: Event Hub Pattern

```
┌─────────┐     POST /send      ┌──────────┐     SSE stream     ┌──────────┐
│  Sender  │ ──────────────────→ │ EventHub │ ──────────────────→ │ Client 1 │
│ (curl)   │                     │ (broker) │ ──────────────────→ │ Client 2 │
└─────────┘                     │          │ ──────────────────→ │ Client 3 │
                                 └──────────┘                     └──────────┘
                                      ↑
                                 Subscribe / Unsubscribe
                                 per client connection
```

- **EventHub** — A broker that manages pub/sub between senders and SSE clients
- **POST /send** — REST endpoint for triggering events
- **GET /events** — SSE endpoint that clients connect to for receiving events

---

## 🛠️ Step-by-Step Guide

### Step 1: Open the starter project

```bash
cd 04-advanced/04-sse/starter
```

Open the `main.go` file — you'll see a structure with `// TODO:` comments indicating where you need to add code.

### Step 2: Understand the Event Hub

The starter includes an `EventHub` struct with method signatures already prepared:

- `Subscribe()` — Register a new client, returns a channel
- `Unsubscribe()` — Remove a client, close the channel
- `Broadcast()` — Send an event to all connected clients
- `ClientCount()` — Count the number of connected clients

### Step 3: Fill in the Broadcast logic

Replace the `// TODO:` in `Broadcast()`:

```go
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
```

> 🔑 Use `select` with `default` for non-blocking send — if a client's buffer is full, it will be skipped instead of blocking

### Step 4: Write the SSE Stream Handler

This is the heart of the lesson! Replace the `// TODO:` in the SSE stream endpoint:

```go
app.Get("/events", func(c *kruda.Ctx) error {
    ch := hub.Subscribe()
    defer hub.Unsubscribe(ch)

    return c.SSE(func(stream *kruda.SSEStream) error {
        // Send a welcome event
        stream.Event("connected", map[string]any{
            "message": "connected",
            "clients": hub.ClientCount(),
        })

        // Loop waiting for events until the client disconnects
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
```

> 📡 `c.SSE()` takes a callback function that receives `*kruda.SSEStream` — inside the callback you can send events and wait for client disconnect via `stream.Done()`

> 🔄 `stream.Done()` fires when the client disconnects — use it to clean up the connection

### Step 5: Write the Send Event Handler

Replace the `// TODO:` for POST /send:

```go
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
```

### Step 6: Write the Client Count Handler and Heartbeat

Replace the `// TODO:` for GET /clients and add the heartbeat goroutine:

```go
kruda.Get[struct{}, MessageResponse](app, "/clients",
    func(c *kruda.C[struct{}]) (*MessageResponse, error) {
        return &MessageResponse{
            Message: fmt.Sprintf("%d client(s) connected", hub.ClientCount()),
        }, nil
    },
)

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
```

> 💓 A heartbeat every 15 seconds keeps the connection active — many reverse proxies will close idle connections

### Step 7: Run and test

```bash
go run main.go
```

Open terminal window 1 — connect to the SSE stream:

```bash
curl -N http://localhost:3000/events
```

Open terminal window 2 — send events:

```bash
# Send an event with a custom name
curl -X POST http://localhost:3000/send \
  -H "Content-Type: application/json" \
  -d '{"event":"notification","data":"{\"title\":\"Hello SSE!\"}"}'

# Send an event with the default name (message)
curl -X POST http://localhost:3000/send \
  -H "Content-Type: application/json" \
  -d '{"data":"simple message"}'

# Check the number of connected clients
curl http://localhost:3000/clients
```

In terminal window 1 you'll see events appearing in real-time:

```
event: connected
data: {"message":"connected","clients":1}

event: notification
data: {"title":"Hello SSE!"}
id: 1

event: message
data: simple message
id: 2
```

🎉 Congratulations! You've built an SSE server that pushes events in real-time!

### Test with Browser (EventSource API)

Open the browser console and run:

```javascript
const es = new EventSource("http://localhost:3000/events");

es.addEventListener("connected", (e) => {
  console.log("Connected:", JSON.parse(e.data));
});

es.addEventListener("notification", (e) => {
  console.log("Notification:", e.data);
});

es.addEventListener("message", (e) => {
  console.log("Message:", e.data);
});

es.addEventListener("heartbeat", (e) => {
  console.log("Heartbeat:", JSON.parse(e.data));
});
```

---

## 🔍 Compare with complete/

If you get stuck anywhere, check the reference implementation in **[complete/](./complete/)** and compare it with your code, or use the `diff` command:

```bash
diff starter/main.go complete/main.go
```

---

## 💡 Key Takeaways

| Concept | Description |
|---|---|
| `c.SSE(func(stream *kruda.SSEStream) error {...})` | Opens an SSE stream with a callback pattern and sets the necessary headers |
| `stream.Event(name, data)` | Sends an event to the client in SSE format |
| `stream.EventWithID(id, name, data)` | Sends an event with an ID for reconnection |
| `stream.Done()` | Detects when the client disconnects (used inside the SSE callback) |
| Event Hub pattern | Pub/sub broker for fan-out events to multiple clients |
| Non-blocking send | Uses `select` + `default` to prevent a slow client from blocking others |
| Heartbeat | Sends periodic events to keep the connection active |
| `EventSource` API | Browser API for receiving SSE events with auto-reconnect |

---

## ➡️ Next Lesson

Awesome! 🎊 You've learned how to push real-time data with Kruda's SSE. In the next lesson we'll learn about the **MCP Server** — how to build a Model Context Protocol server with Kruda 🤖

👉 [Section 04-05 — MCP Server](../05-mcp-server/)
