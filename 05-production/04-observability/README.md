# 🔭 Section 05-04 — Observability: Turnkey Tracing, Metrics & Health with contrib/observability

⏱️ Estimated time: **30 minutes**

Welcome to the final lesson of the Tutorial! In this lesson you will learn how to wire OpenTelemetry **tracing**, Prometheus-style **metrics**, and Kubernetes **health probes** into a Kruda app with a single call using `contrib/observability` -- a turnkey alternative to Section 05-01's manual `contrib/prometheus` setup.

---

## Target Learning Outcomes

- Understand what `contrib/observability` provides beyond `contrib/prometheus` (Section 05-01)
- Call `observability.Enable(app, cfg)` **before** registering routes
- Expose Prometheus metrics at `/metrics` without wiring middleware by hand
- Expose Kubernetes liveness/readiness probes at `/livez` and `/readyz`
- View real distributed-tracing spans on stdout with zero external infrastructure

---

## What You Will Learn

By the end of this lesson you will be able to:

- Wire tracing, metrics, and health probes in one call: `observability.Enable(app, cfg)`
- Explain why `Enable` must run before any route registration
- Configure `observability.Config` (`ServiceName`, `MetricsPublic`, etc.)
- Flush pending trace data on shutdown with `prov.Flush(ctx)`
- View a real trace span on stdout using `OTEL_TRACES_EXPORTER=console`
- Recognise a known upstream quirk in the module's route-label reporting

---

## Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25+ |
| Git | Latest |
| Text Editor / IDE | VS Code, GoLand, or your preferred editor |
| Terminal | For running `go run` and setting environment variables |

> If you haven't completed Section 05-03 (Benchmark), consider going back first -- see [Section 05-03 -- Benchmark](../03-benchmark/)

---

## File Structure

```
05-production/04-observability/
├── README.md              <-- You are here
├── starter/               <-- Starter code (with TODOs to fill in)
│   ├── go.mod
│   └── main.go
└── complete/              <-- Complete solution
    ├── go.mod
    └── main.go
```

- **[starter/](./starter/)** -- Skeleton code that compiles but has `// TODO:` markers for you to complete
- **[complete/](./complete/)** -- Full working solution

---

## Why contrib/observability?

Section 05-01 wired Prometheus metrics by hand with `contrib/prometheus` -- one middleware, one `/metrics` handler. `contrib/observability` is a heavier, turnkey alternative that adds **distributed tracing** and **Kubernetes health probes** on top, in a single call:

| | `contrib/prometheus` (05-01) | `contrib/observability` (this lesson) |
|---|---|---|
| Metrics | ✅ manual middleware + handler | ✅ automatic |
| Tracing | ❌ | ✅ OpenTelemetry spans |
| Health probes | ❌ (you write your own) | ✅ `/livez`, `/readyz` built in |
| Setup | `app.Use(prometheus.New())` + `app.Get("/metrics", ...)` | one `observability.Enable(app, cfg)` call |
| Dependency weight | Light (~20 modules) | Heavier (~40 modules -- full OTel SDK) |

> Use `contrib/prometheus` when you only need metrics and want a light dependency footprint. Use `contrib/observability` when you also want tracing and Kubernetes probes out of the box.

---

## Step-by-Step Guide

### Step 1: Open the starter project

```bash
cd 05-production/04-observability/starter
```

Open `main.go` -- you will see the skeleton with a `// TODO:` comment for enabling observability.

### Step 2: Enable Observability Before Registering Routes

```go
app := kruda.New()

prov, err := observability.Enable(app, observability.Config{
    ServiceName:   "kruda-tutorial-observability",
    MetricsPublic: true, // mount /metrics on this app's own port
})
if err != nil {
    log.Fatalf("observability.Enable: %v", err)
}
defer prov.Flush(context.Background())
```

> ⚠️ `observability.Enable` **must** run before any route is registered. It installs its span-tracking middleware via `app.Use(...)`, and kruda bakes middleware into a route's handler chain at registration time -- a route registered before `Enable` runs without instrumentation.

> `MetricsPublic: true` mounts `/metrics` on the app's own port so `curl localhost:3000/metrics` works for this lesson. The default (`false`) instead starts a separate internal loopback-only listener -- appropriate in production where metrics scraping happens on a private network, not the public port.

### Step 3: Register a Route

```go
kruda.Get[struct{}, MessageResponse](app, "/hello", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
    return &MessageResponse{
        Message: fmt.Sprintf("hello at %s", time.Now().UTC().Format(time.RFC3339)),
    }, nil
})
```

Because this is registered *after* `Enable`, every request to `/hello` gets a trace span and is counted in the automatic RED (Rate, Errors, Duration) metrics.

### Step 4: Run and Test Metrics + Health

```bash
go run main.go
```

```bash
# Your route
curl http://localhost:3000/hello

# Prometheus metrics
curl http://localhost:3000/metrics

# Kubernetes probes
curl http://localhost:3000/livez
curl http://localhost:3000/readyz
```

`/livez` and `/readyz` return immediately with no dependencies:

```json
{"status":"ok"}
```

```json
{"status":"ok","checks":{}}
```

`/metrics` returns real Prometheus exposition text, including `http_server_active_requests` and `http_server_request_duration_seconds`.

> ⚠️ **Known upstream quirk:** at the time of writing, `contrib/observability` v1.0.0 truncates the `http.route` label for short paths -- e.g. `/metrics` may appear in its own metrics as `http_route="/metri"`. This is a bug in the module itself (its route-templating logic), not something you need to fix in your code. It doesn't affect the `/hello` route's longer path in practice, but don't be alarmed if you see a truncated label on a short one.

### Step 5: See a Real Trace Span

By default, tracing exports over OTLP -- with no collector configured, spans are silently dropped (with a single startup warning in the log). To see tracing work with zero external infrastructure, export to stdout instead:

```bash
OTEL_TRACES_EXPORTER=console go run main.go
```

```bash
curl http://localhost:3000/hello
```

You will see a JSON span printed to stdout, including `trace_id`, `span_id`, request attributes, and the `service.name` you configured -- no Jaeger, Tempo, or OTel Collector required.

---

## Compare with complete/

If you get stuck, check the solution in **[complete/](./complete/)** and compare:

```bash
diff starter/main.go complete/main.go
```

---

## Key Concepts Summary

| Concept | Description |
|---|---|
| `observability.Enable(app, cfg)` | Wires tracing, RED metrics, and health probes in one call |
| Register `Enable` before routes | Middleware bakes into the chain at registration time |
| `observability.Config` | `ServiceName`, `MetricsPublic`, `MetricsPath`, `LivenessPath`, `ReadinessPath`, and more |
| `prov.Flush(ctx)` | Bounded shutdown flush for pending trace/metric data |
| `/livez` / `/readyz` | Kubernetes liveness/readiness probe endpoints, zero dependencies |
| `/metrics` | Prometheus exposition endpoint, automatic RED metrics |
| `OTEL_TRACES_EXPORTER=console` | Print trace spans to stdout with no external collector |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Where traces export by default -- unset means spans are silently dropped |

---

## Congratulations! You have completed the Tutorial!

You have learned everything in the Kruda Tutorial from beginner to production!

What you learned across the tutorial:

- **Beginner** -- REST API + Typed Handlers
- **Auto CRUD** -- Generate CRUD endpoints automatically
- **Intermediate** -- Database, Config, Error Handling
- **Advanced** -- DI Container, Auth, OpenAPI, SSE, MCP, WebSocket, Testing, Architecture
- **Production** -- Monitoring, Docker Deploy, Benchmark, Observability

Happy building with Kruda!

--> [Back to main page](../../)
