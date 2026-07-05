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

	// TODO: Enable observability BEFORE registering any route.
	//
	// observability.Enable MUST run before route registration --
	// it installs span middleware via app.Use, and kruda bakes
	// middleware into a route's chain at registration time.
	//
	// Example:
	//   prov, err := observability.Enable(app, observability.Config{
	//       ServiceName:   "kruda-tutorial-observability",
	//       MetricsPublic: true, // mount /metrics on this app's own port
	//   })
	//   if err != nil {
	//       log.Fatalf("observability.Enable: %v", err)
	//   }
	//   defer prov.Flush(context.Background())
	_ = context.Background
	_ = observability.Enable

	// TODO: Register GET /hello -- returns a MessageResponse.
	//
	// Example:
	//   kruda.Get[struct{}, MessageResponse](app, "/hello", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
	//       return &MessageResponse{
	//           Message: fmt.Sprintf("hello at %s", time.Now().UTC().Format(time.RFC3339)),
	//       }, nil
	//   })
	_ = fmt.Sprintf
	_ = time.Now

	log.Println("Server starting on :3000 ...")
	log.Fatal(app.Listen(":3000"))
}
