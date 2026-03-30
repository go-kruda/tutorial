package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Request / Response Types
// ============================================================
//
// These types define the JSON shape of our API responses.
// When running inside a Docker container, the application
// behaves identically to running on bare metal -- Kruda's
// Wing Transport works seamlessly in containerised environments.

// HealthResponse represents the JSON response for the health
// check endpoint. Container orchestrators (Docker Compose,
// Kubernetes) use health checks to determine if a container
// is ready to receive traffic.
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// MessageResponse is a generic message envelope.
type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// Helpers
// ============================================================

// getEnv reads an environment variable with a fallback default.
// This pattern is essential for containerised applications where
// configuration is injected via environment variables (12-factor
// app methodology).
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	// -- 1. Read Configuration from Environment ------------
	//
	// In a Docker container, configuration is passed via
	// environment variables (docker run -e PORT=8080 ...).
	// This follows the 12-factor app methodology:
	// https://12factor.net/config
	port := getEnv("PORT", "3000")

	// -- 2. Create the Kruda Application -------------------
	//
	// kruda.New() initialises the app with Wing Transport --
	// Kruda's high-performance epoll-based networking layer.
	// This works identically inside a Docker container.
	app := kruda.New()

	// -- 3. Register Routes --------------------------------
	//
	// The /health endpoint is critical for Docker HEALTHCHECK
	// and orchestrator probes. The /hello endpoint demonstrates
	// a simple application route.

	// Health check endpoint -- container orchestrators
	// (Docker, Kubernetes) periodically probe this endpoint
	// to verify the application is alive and ready. If the
	// health check fails, the orchestrator can restart the
	// container automatically.
	kruda.Get[struct{}, HealthResponse](app, "/health", func(c *kruda.C[struct{}]) (*HealthResponse, error) {
		return &HealthResponse{
			Status:    "ok",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Version:   getEnv("APP_VERSION", "1.0.0"),
		}, nil
	})

	// Greeting endpoint -- includes the hostname, which is
	// useful for verifying which container is responding when
	// running multiple replicas behind a load balancer.
	kruda.Get[struct{}, MessageResponse](app, "/hello", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
		hostname, _ := os.Hostname()
		return &MessageResponse{
			Message: fmt.Sprintf("Hello from Kruda running in Docker! (host: %s)", hostname),
		}, nil
	})

	// -- 4. Start the Server -------------------------------
	//
	// Listen on 0.0.0.0 (all interfaces) so the container's
	// port mapping works correctly. Docker maps the container
	// port to the host via -p flag (e.g., -p 8080:3000).
	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Running in Docker -- health check at http://localhost:%s/health\n", port)
	log.Printf("Server starting on %s ...", addr)
	log.Fatal(app.Listen(addr))
}
