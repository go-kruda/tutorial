package main

import (
	"log"
	"sync"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Request / Response Types
// ============================================================
//
// These types are used by both the application and the benchmark
// tests. Keeping them in the same package lets benchmark_test.go
// access them directly -- no need for a separate test helper.

// CreateItemInput represents the JSON body for creating an item.
type CreateItemInput struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// ItemResponse represents the JSON response returned for an item.
type ItemResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// MessageResponse is a generic message envelope.
type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// In-Memory Item Store
// ============================================================
//
// A simple thread-safe store used by the handlers. The benchmark
// tests exercise these handlers under load to measure throughput
// and memory allocation.

// ItemStore provides thread-safe in-memory storage for items.
// Benchmarks create a fresh store per iteration to avoid
// cross-contamination between runs.
type ItemStore struct {
	mu     sync.Mutex
	items  []ItemResponse
	nextID int
}

// NewItemStore creates an empty ItemStore.
func NewItemStore() *ItemStore {
	return &ItemStore{nextID: 1}
}

// Create adds a new item and returns it with an assigned ID.
func (s *ItemStore) Create(name string, price int) ItemResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := ItemResponse{
		ID:    s.nextID,
		Name:  name,
		Price: price,
	}
	s.nextID++
	s.items = append(s.items, item)
	return item
}

// All returns a copy of all items in the store.
func (s *ItemStore) All() []ItemResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]ItemResponse, len(s.items))
	copy(result, s.items)
	return result
}

// ============================================================
// Application Factory
// ============================================================
//
// NewApp creates a configured Kruda application with all routes
// registered. This factory pattern allows benchmark tests to
// create a fresh app instance and use kruda.NewTestClient to
// exercise the full HTTP stack without starting a real server.

// NewApp creates and returns a fully configured Kruda app.
//
// Why a factory function?
// ----------------------
// Benchmark tests need to create the same application repeatedly.
// By extracting the setup into NewApp, both main() and the
// benchmarks share the exact same configuration -- ensuring
// benchmark results accurately reflect production behaviour.
func NewApp(store *ItemStore) *kruda.App {
	app := kruda.New()

	// GET /health -- simple health check, useful as a baseline
	// to understand the minimum overhead of a Kruda handler.
	kruda.Get[struct{}, MessageResponse](app, "/health", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
		return &MessageResponse{Message: "ok"}, nil
	})

	// GET /items -- returns all items. The benchmark for this
	// route measures serialisation throughput: how fast Kruda
	// can convert a Go slice to a JSON array response.
	kruda.Get[struct{}, []ItemResponse](app, "/items", func(c *kruda.C[struct{}]) (*[]ItemResponse, error) {
		items := store.All()
		return &items, nil
	})

	// POST /items -- creates a new item. Kruda automatically
	// deserialises the JSON body into CreateItemInput. The
	// benchmark measures the full cycle: JSON parsing, handler
	// logic, and response serialisation.
	kruda.Post[CreateItemInput, ItemResponse](app, "/items", func(c *kruda.C[CreateItemInput]) (*ItemResponse, error) {
		item := store.Create(c.In.Name, c.In.Price)
		return &item, nil
	})

	app.Compile()
	return app
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	store := NewItemStore()

	// Create a new Kruda application with Wing Transport.
	// Wing Transport uses epoll under the hood -- the benchmark
	// numbers you see reflect this high-performance networking.
	app := NewApp(store)

	log.Println("Benchmark app starting on :3000 ...")
	log.Fatal(app.Listen(":3000"))
}
