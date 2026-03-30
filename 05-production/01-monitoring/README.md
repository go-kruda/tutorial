# 📊 Section 05-01 — Monitoring: Prometheus + Grafana

⏱️ Estimated time: **30 minutes**

Welcome to the first Production section! In this lesson you will learn how to add **Prometheus metrics** to a Kruda application using the `contrib/prometheus` package, and create a **Grafana dashboard** to visualise request rate, latency, and active requests in real time.

---

## Learning Objectives

- Understand the concept of **observability** and why monitoring matters for production
- Learn how to add Prometheus metrics middleware with `prometheus.New()`
- Serve a `/metrics` endpoint with `prometheus.Handler()`
- Understand the automatically tracked metrics: **Counter**, **Histogram**, and **Gauge**
- Import a Grafana dashboard JSON to visualise metrics

---

## What You Will Learn

By the end of this lesson you will be able to:

- Add Prometheus metrics middleware with `app.Use(prometheus.New())`
- Serve the `/metrics` endpoint with `app.Get("/metrics", prometheus.Handler())`
- Understand auto-tracked Counter (`http_requests_total`) for request counts
- Understand auto-tracked Histogram (`http_request_duration_seconds`) for latency distribution
- Understand auto-tracked Gauge (`http_requests_in_flight`) for active requests
- Configure Prometheus to scrape metrics from a Kruda app
- Import a Grafana dashboard JSON to view request rate, latency percentiles, and heatmaps

---

## Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25+ |
| Docker & Docker Compose | For running Prometheus + Grafana |
| Git | Latest |
| Text Editor / IDE | VS Code, GoLand, or your preferred editor |

> If you haven't completed Section 04, consider going back first as this lesson builds on Middleware Chain knowledge -- see [Section 04-08 -- Architecture](../../04-advanced/08-architecture/)

---

## File Structure

```
05-production/01-monitoring/
|-- README.md              <-- You are here
|-- starter/               <-- Starter code (with TODOs to fill in)
|   |-- go.mod
|   +-- main.go
+-- complete/              <-- Complete solution
    |-- go.mod
    |-- main.go
    +-- dashboard.json     <-- Grafana dashboard definition
```

- **[starter/](./starter/)** -- Skeleton code that compiles but has `// TODO:` markers for you to complete
- **[complete/](./complete/)** -- Full working solution with Grafana dashboard JSON

---

## Why Monitor?

When your application is in production, you need to answer these questions:

| Question | Metric That Answers It |
|---|---|
| How many requests per second? | Counter (`http_requests_total`) |
| What does response time look like? | Histogram (`http_request_duration_seconds`) |
| How many requests are in flight right now? | Gauge (`http_requests_in_flight`) |
| Is p99 latency exceeding SLA? | Histogram quantile (p99) |

> **Prometheus** stores time-series data, **Grafana** displays it as beautiful dashboards -- both are open-source industry standards

---

## Step-by-Step Guide

### Step 1: Open the starter project

```bash
cd 05-production/01-monitoring/starter
```

Open `main.go` -- you will see the skeleton with `// TODO:` comments.

### Step 2: Add Prometheus Metrics Middleware

The `contrib/prometheus` package provides a middleware that automatically instruments every request:

```go
import (
    "github.com/go-kruda/kruda"
    "github.com/go-kruda/kruda/contrib/prometheus"
)

app := kruda.New()

// Add Prometheus middleware -- automatically tracks:
// - http_requests_total (Counter)
// - http_request_duration_seconds (Histogram)
// - http_requests_in_flight (Gauge)
// - http_response_size_bytes (Histogram)
app.Use(prometheus.New())
```

Optional custom configuration:

```go
app.Use(prometheus.New(prometheus.Config{
    Namespace: "myapp",
    Subsystem: "http",
}))
```

> The middleware runs for EVERY request automatically -- no per-handler changes needed

### Step 3: Serve the Metrics Endpoint

```go
app.Get("/metrics", prometheus.Handler())
```

> `prometheus.Handler()` returns a handler that serves metrics in Prometheus exposition format

### Step 4: Register Application Routes

```go
kruda.Get[struct{}, HealthResponse](app, "/health", func(c *kruda.C[struct{}]) (*HealthResponse, error) {
    return &HealthResponse{
        Status:    "ok",
        Timestamp: time.Now().UTC().Format(time.RFC3339),
    }, nil
})

kruda.Get[struct{}, MessageResponse](app, "/hello", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
    return &MessageResponse{
        Message: "Hello from Kruda with Prometheus monitoring!",
    }, nil
})
```

### Step 5: Start the Server

```go
log.Fatal(app.Listen(":3000"))
```

### Step 6: Run and Test

```bash
go run main.go
```

Open another terminal and test:

```bash
# Send requests to the application
curl http://localhost:3000/hello
curl http://localhost:3000/health

# View raw Prometheus metrics
curl http://localhost:3000/metrics
```

You will see output like this from `/metrics`:

```
# HELP http_requests_total Total number of HTTP requests.
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/hello",status="200"} 1

# HELP http_request_duration_seconds HTTP request latency.
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{le="0.001"} 1
...
```

### Step 7: Set Up Prometheus (Optional)

Create a `prometheus.yml` file:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: "kruda-app"
    static_configs:
      - targets: ["host.docker.internal:3000"]
```

Run Prometheus with Docker:

```bash
docker run -d --name prometheus \
  -p 9090:9090 \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus
```

### Step 8: Import Grafana Dashboard

1. Run Grafana with Docker:

```bash
docker run -d --name grafana \
  -p 3001:3000 \
  grafana/grafana
```

2. Open http://localhost:3001 (login: admin/admin)
3. Add Prometheus data source: http://host.docker.internal:9090
4. Import dashboard from `complete/dashboard.json`

You will see a dashboard showing:
- **Request Rate** -- requests per second
- **Latency Percentiles** -- p50, p90, p99
- **Latency Heatmap** -- response time distribution
- **Active Requests** -- currently processing
- **Total Requests** -- cumulative count

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
| `prometheus.New()` | Middleware that automatically instruments all HTTP requests |
| `prometheus.Handler()` | Serves `/metrics` endpoint in Prometheus exposition format |
| `http_requests_total` | Auto-tracked Counter -- total requests by method/path/status |
| `http_request_duration_seconds` | Auto-tracked Histogram -- latency distribution |
| `http_requests_in_flight` | Auto-tracked Gauge -- currently processing requests |
| `http_response_size_bytes` | Auto-tracked Histogram -- response body sizes |
| `prometheus.Config` | Optional config for namespace/subsystem prefixes |
| `histogram_quantile()` | PromQL function for computing percentiles from histograms |

---

## Next Lesson

Great work! You have learned how to add Prometheus monitoring to a Kruda app and create a Grafana dashboard. In the next lesson you will learn **Docker Deploy** -- how to containerise a Kruda application with a multi-stage Dockerfile for production deployment.

--> [Section 05-02 -- Docker Deploy](../02-docker-deploy/)
