package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-kruda/kruda"
	"github.com/go-kruda/kruda/contrib/observability"
)

// ============================================================
// Request / Response Types
// ============================================================

// MessageResponse is a generic message envelope.
type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	app := kruda.New()

	// observability.Enable MUST be called before any route is
	// registered -- it installs span middleware via app.Use, and
	// kruda bakes middleware into a route's chain at registration
	// time. Routes registered before Enable would not be
	// instrumented.
	prov, err := observability.Enable(app, observability.Config{
		ServiceName: "kruda-tutorial-observability",
		// MetricsPublic mounts /metrics on this app's own port.
		// The default (false) starts an internal loopback-only
		// listener instead, which is not reachable from outside
		// the process -- fine for a real deployment behind a
		// separate scrape network, but this lesson wants curl to
		// reach it directly.
		MetricsPublic: true,
	})
	if err != nil {
		log.Fatalf("observability.Enable: %v", err)
	}
	// Flush is bounded by Config.FlushTimeout (default 5s) so
	// shutdown never hangs waiting on a stuck exporter.
	defer prov.Flush(context.Background())

	// Registered AFTER Enable, so this route is instrumented --
	// every request gets a trace span and counts toward the RED
	// metrics Enable wires up automatically.
	kruda.Get[struct{}, MessageResponse](app, "/hello", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
		return &MessageResponse{
			Message: fmt.Sprintf("hello at %s", time.Now().UTC().Format(time.RFC3339)),
		}, nil
	})

	log.Println("Server starting on :3000 ...")
	log.Println("  Hello:     GET  /hello")
	log.Println("  Metrics:   GET  /metrics   (Prometheus exposition format)")
	log.Println("  Liveness:  GET  /livez")
	log.Println("  Readiness: GET  /readyz")
	log.Println("  Set OTEL_TRACES_EXPORTER=console to print trace spans to stdout")
	log.Fatal(app.Listen(":3000"))
}
