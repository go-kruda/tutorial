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

// HealthResponse represents the JSON response for the health
// check endpoint. Container orchestrators use this to verify
// the application is alive and ready to receive traffic.
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
// This follows the 12-factor app methodology where configuration
// is injected via environment variables in containers.
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
	// TODO: Read the PORT from environment variable with a
	// default of "3000".
	//
	// Example:
	//   port := getEnv("PORT", "3000")
	_ = fmt.Sprintf // keep fmt import used
	_ = getEnv      // placeholder -- remove after implementing
	_ = time.Now    // placeholder -- remove after implementing

	// TODO: Create a new Kruda application instance.
	//
	// Example:
	//   app := kruda.New()
	var app *kruda.App

	// TODO: Register routes using kruda.Get:
	//   - GET /health -- returns HealthResponse with status, timestamp, version
	//   - GET /hello  -- returns MessageResponse with hostname
	//
	// The /health endpoint is essential for Docker HEALTHCHECK.
	//
	// Example:
	//   kruda.Get[struct{}, HealthResponse](app, "/health", func(c *kruda.C[struct{}]) (*HealthResponse, error) {
	//       return &HealthResponse{
	//           Status:    "ok",
	//           Timestamp: time.Now().UTC().Format(time.RFC3339),
	//           Version:   getEnv("APP_VERSION", "1.0.0"),
	//       }, nil
	//   })
	//
	//   kruda.Get[struct{}, MessageResponse](app, "/hello", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
	//       hostname, _ := os.Hostname()
	//       return &MessageResponse{
	//           Message: fmt.Sprintf("Hello from Docker! (host: %s)", hostname),
	//       }, nil
	//   })

	// Placeholders -- remove after implementing.
	_ = app
	_ = os.Hostname

	// TODO: Start the server on the configured port.
	//
	// Example:
	//   addr := fmt.Sprintf(":%s", port)
	//   log.Fatal(app.Listen(addr))
	log.Println("Server starting on :3000 ...")
}
