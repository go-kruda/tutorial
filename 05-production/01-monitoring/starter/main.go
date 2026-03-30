package main

import (
	"fmt"
	"log"

	"github.com/go-kruda/kruda"
	"github.com/go-kruda/kruda/contrib/prometheus"
)

// ============================================================
// Request / Response Types
// ============================================================

// HealthResponse represents the JSON response for the health check.
type HealthResponse struct {
	Status string `json:"status"`
}

// MessageResponse is a generic message envelope.
type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// Prometheus Monitoring
// ============================================================
//
// Kruda provides a contrib/prometheus package that integrates
// Prometheus monitoring as middleware. It automatically tracks:
//   - http_requests_total          (Counter)
//   - http_request_duration_seconds (Histogram)
//   - http_requests_in_flight       (Gauge)
//   - http_response_size_bytes      (Histogram)

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	// TODO: Create a new Kruda application instance.
	//
	// Example:
	//   app := kruda.New()
	_ = fmt.Sprintf // keep fmt import used
	var app *kruda.App

	// TODO: Add Prometheus metrics middleware.
	// This automatically instruments all requests with metrics.
	//
	// Example:
	//   app.Use(prometheus.New())
	//
	// Optional custom config:
	//   app.Use(prometheus.New(prometheus.Config{
	//       Namespace: "myapp",
	//       Subsystem: "http",
	//   }))

	// TODO: Serve the Prometheus metrics endpoint.
	// prometheus.Handler() returns a handler that serves metrics
	// in Prometheus exposition format.
	//
	// Example:
	//   app.Get("/metrics", prometheus.Handler())

	// TODO: Register application routes using kruda.Get.
	//
	// Example:
	//   kruda.Get[struct{}, HealthResponse](app, "/health", func(c *kruda.C[struct{}]) (*HealthResponse, error) {
	//       return &HealthResponse{Status: "ok"}, nil
	//   })
	//
	//   kruda.Get[struct{}, MessageResponse](app, "/hello", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
	//       return &MessageResponse{Message: "Hello with monitoring!"}, nil
	//   })

	// Placeholders -- remove after implementing.
	_ = app
	_ = prometheus.New
	_ = prometheus.Handler

	// TODO: Start the server on port 3000.
	//
	// Example:
	//   log.Fatal(app.Listen(":3000"))
	log.Println("Server starting on :3000 ...")
	log.Println("Metrics available at http://localhost:3000/metrics")
}
