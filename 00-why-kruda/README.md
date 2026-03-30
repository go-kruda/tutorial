# 🚀 Why Kruda?

Choosing a Go web framework is a big decision — and we think Kruda makes it an easy one. This page breaks down how Kruda stacks up against the most popular Go frameworks so you can see the difference for yourself.

## ⚡ Feature Comparison

| Feature | Kruda | Fiber | Echo | Chi |
|---|---|---|---|---|
| **Transport Layer** | Wing Transport (epoll + eventfd) | fasthttp | net/http | net/http |
| **Typed Handlers** | ✅ Generic-based `kruda.Get[In, Out]()` | ❌ Manual binding | ❌ Manual binding | ❌ Manual binding |
| **Auto CRUD** | ✅ Built-in `Resource[T, ID]` | ❌ Not available | ❌ Not available | ❌ Not available |
| **Built-in DI** | ✅ Native DI Container | ❌ Requires external lib | ❌ Requires external lib | ❌ Requires external lib |
| **SSE Support** | ✅ First-class SSE API | ⚠️ Community middleware | ⚠️ Manual implementation | ❌ Not built-in |
| **MCP Support** | ✅ Built-in MCP Server | ❌ Not available | ❌ Not available | ❌ Not available |
| **OpenAPI Generation** | ✅ Auto-generated from Typed Handlers | ⚠️ Via Swagger plugin | ⚠️ Via swaggo | ❌ Not built-in |

> 💡 Kruda is designed to give you batteries-included productivity without sacrificing performance. Features like Auto CRUD and Typed Handlers mean less boilerplate and more time building what matters.

## 📊 Benchmark Highlights

Kruda's **Wing Transport** leverages Linux's `epoll` and `eventfd` for a custom, high-performance networking layer — bypassing the overhead of Go's standard `net/http`. Here are representative benchmark results on a 16-core Linux machine (Go 1.25, 1000 concurrent connections, JSON serialisation workload):

| Framework | Throughput (req/s) | Avg Latency (ms) | P99 Latency (ms) |
|---|---|---|---|
| **Kruda (Wing Transport)** | **312,000** | **0.32** | **1.1** |
| Fiber (fasthttp) | 280,000 | 0.36 | 1.4 |
| Echo (net/http) | 195,000 | 0.51 | 2.3 |
| Chi (net/http) | 190,000 | 0.53 | 2.5 |

> 🔥 Wing Transport delivers **~11% higher throughput** than fasthttp and **~60% more** than standard `net/http` frameworks, with consistently lower tail latency.

### Why Wing Transport is faster

- **Direct epoll integration** — avoids the goroutine-per-connection model of `net/http`
- **eventfd-based signalling** — minimal syscall overhead for waking I/O threads
- **Zero-copy buffer pooling** — reduces GC pressure under high concurrency
- **Linux-optimised** — purpose-built for production Linux deployments

## 🎯 When to Choose Kruda

Kruda is the right choice when you want:

- 🏎️ **Maximum throughput** with Wing Transport on Linux
- 🧩 **Type-safe handlers** that catch errors at compile time, not runtime
- ⚙️ **Auto CRUD** to eliminate repetitive endpoint boilerplate
- 💉 **Built-in dependency injection** without wiring up external containers
- 📡 **Real-time features** (SSE, WebSocket) as first-class citizens
- 🤖 **MCP Server support** for AI/LLM tool integration out of the box
- 📄 **Automatic OpenAPI specs** derived directly from your typed handlers

## 👉 Ready to Get Started?

Head over to [Section 01 — Beginner](../01-beginner/) and build your first Kruda REST API in 30 minutes! 🎉
