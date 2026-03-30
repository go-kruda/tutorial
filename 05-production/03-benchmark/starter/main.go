package main

import (
	"log"
	"sync"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Request / Response Types
// ============================================================

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

// ItemStore provides thread-safe in-memory storage for items.
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

// NewApp creates and returns a fully configured Kruda app.
// This factory pattern allows benchmark tests to create a fresh
// app instance and use kruda.NewTestClient to test the full
// HTTP stack without starting a real server.
func NewApp(store *ItemStore) *kruda.App {
	app := kruda.New()

	// TODO: Register routes using kruda.Get and kruda.Post:
	//
	//   GET  /health -- returns MessageResponse with Message "ok"
	//   GET  /items  -- returns all items from store.All()
	//   POST /items  -- creates item using c.In.Name and c.In.Price
	//
	// Example:
	//   kruda.Get[struct{}, MessageResponse](app, "/health", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
	//       return &MessageResponse{Message: "ok"}, nil
	//   })
	//
	//   kruda.Get[struct{}, []ItemResponse](app, "/items", func(c *kruda.C[struct{}]) (*[]ItemResponse, error) {
	//       items := store.All()
	//       return &items, nil
	//   })
	//
	//   kruda.Post[CreateItemInput, ItemResponse](app, "/items", func(c *kruda.C[CreateItemInput]) (*ItemResponse, error) {
	//       item := store.Create(c.In.Name, c.In.Price)
	//       return &item, nil
	//   })

	// Placeholder -- remove after implementing.
	_ = store

	// TODO: Call app.Compile() after registering all routes.
	// This is required for kruda.NewTestClient to work in benchmarks.
	//
	// Example:
	//   app.Compile()

	return app
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	// TODO: Create an ItemStore and build the app using NewApp.
	//
	// Example:
	//   store := NewItemStore()
	//   app := NewApp(store)
	//   log.Fatal(app.Listen(":3000"))

	// Placeholders -- remove after implementing.
	_ = NewItemStore

	log.Println("Benchmark app starting on :3000 ...")
}
