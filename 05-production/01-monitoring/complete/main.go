package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-kruda/kruda"
	"github.com/go-kruda/kruda/contrib/prometheus"
)

// ============================================================
// Request / Response Types
// ============================================================
//
// Simple types for our demo endpoints. In a real application
// these would be your domain models -- the monitoring layer
// wraps around them transparently.

// HealthResponse represents the JSON response for the health check.
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// MessageResponse is a generic message envelope.
type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// Prometheus Monitoring with contrib/prometheus
// ============================================================
//
// Kruda provides a contrib/prometheus package that integrates
// Prometheus monitoring as middleware. When you add the
// prometheus.New() middleware, it automatically tracks:
//
//   - http_requests_total          (Counter)
//   - http_request_duration_seconds (Histogram)
//   - http_requests_in_flight       (Gauge)
//   - http_response_size_bytes      (Histogram)
//
// Why use the contrib package instead of manual metrics?
// -----------------------------------------------------
// The contrib/prometheus package provides battle-tested,
// standardised metrics out of the box. You get consistent
// metric names that follow Prometheus naming conventions,
// automatic label injection (method, path, status code),
// and a ready-made /metrics endpoint handler.
//
// For custom metrics beyond HTTP request tracking, you can
// still use the prometheus client library directly.

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	// -- 1. Create the Kruda Application -------------------
	app := kruda.New()

	// -- 2. Add Prometheus Metrics Middleware ---------------
	//
	// prometheus.New() returns a middleware that automatically
	// instruments every HTTP request with Prometheus metrics:
	//
	//   - http_requests_total: Counter of total requests,
	//     labelled by method, path, and status code.
	//   - http_request_duration_seconds: Histogram of request
	//     latency with default bucket boundaries.
	//   - http_requests_in_flight: Gauge tracking how many
	//     requests are currently being processed.
	//   - http_response_size_bytes: Histogram of response sizes.
	//
	// You can optionally customise the namespace and subsystem:
	//
	//   app.Use(prometheus.New(prometheus.Config{
	//       Namespace: "myapp",
	//       Subsystem: "http",
	//   }))
	//
	// This would prefix all metric names with "myapp_http_".
	app.Use(prometheus.New())

	// -- 3. Serve Prometheus Metrics Endpoint ---------------
	//
	// prometheus.Handler() returns a handler that serves metrics
	// in Prometheus exposition format. Point your Prometheus
	// scraper at this endpoint.
	//
	// After this call, visiting http://localhost:3000/metrics
	// returns output like:
	//
	//   # HELP http_requests_total Total number of HTTP requests.
	//   # TYPE http_requests_total counter
	//   http_requests_total{method="GET",path="/hello",status="200"} 42
	//
	//   # HELP http_request_duration_seconds HTTP request latency.
	//   # TYPE http_request_duration_seconds histogram
	//   http_request_duration_seconds_bucket{le="0.001"} 10
	//   ...
	//
	// Prometheus scrapes this endpoint at a configured interval
	// (typically every 15s) and stores the time-series data.
	app.Get("/metrics", prometheus.Handler())

	// -- 4. Register Application Routes --------------------
	//
	// These are your normal application endpoints. The
	// prometheus middleware wraps them transparently --
	// every request is automatically instrumented.

	// Health check endpoint -- essential for:
	//   - Load balancer health checks
	//   - Kubernetes liveness/readiness probes
	//   - Uptime monitoring services
	kruda.Get[struct{}, HealthResponse](app, "/health", func(c *kruda.C[struct{}]) (*HealthResponse, error) {
		return &HealthResponse{
			Status:    "ok",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}, nil
	})

	// Simple greeting endpoint -- demonstrates that metrics
	// are collected automatically for ALL routes via the
	// middleware.
	kruda.Get[struct{}, MessageResponse](app, "/hello", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
		return &MessageResponse{
			Message: "Hello from Kruda with Prometheus monitoring!",
		}, nil
	})

	// -- 5. Start the Server -------------------------------
	//
	// Once the server is running, you can:
	//   1. Hit application endpoints:
	//      curl http://localhost:3000/hello
	//      curl http://localhost:3000/health
	//
	//   2. View raw Prometheus metrics:
	//      curl http://localhost:3000/metrics
	//
	//   3. Configure Prometheus to scrape localhost:3000/metrics
	//
	//   4. Import the provided Grafana dashboard JSON to
	//      visualise request rates, latency percentiles, and
	//      active request counts in real time.
	fmt.Println("Metrics available at http://localhost:3000/metrics")
	log.Println("Server starting on :3000 ...")
	log.Fatal(app.Listen(":3000"))
}
