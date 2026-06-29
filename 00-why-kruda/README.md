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

Kruda's **Wing Transport** leverages Linux's `epoll` and `eventfd` for a custom, high-performance networking layer — bypassing the overhead of Go's standard `net/http`. Here are representative benchmark results on a 16-core Linux machine (Go 1.25.11, CPU-bound handler workloads):

| Workload | Kruda (Wing) | Fiber (fasthttp) | Actix (Rust) |
|---|---|---|---|
| **Plaintext** | **846K req/s** | 670K req/s | 814K req/s |
| **JSON** | **805K req/s** | 625K req/s | 790K req/s |
| **DB (read)** | **108K req/s** | 107K req/s | 37K req/s |

> 🔥 Wing Transport delivers **~12% higher throughput** than Actix (Rust) and **~29% more** than Fiber on JSON workloads, with consistently lower p99 latency.

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
