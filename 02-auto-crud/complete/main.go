package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Model Definition
// ============================================================
//
// Auto CRUD works by registering a ResourceService implementation
// for your model. Kruda generates all 5 CRUD endpoints from the
// service interface -- you define how data is stored, and the
// framework handles HTTP routing, JSON serialisation, pagination,
// and error responses automatically.
//
// Why Auto CRUD?
// --------------
// In a traditional REST API you write the same boilerplate for
// every resource: list, create, get-by-id, update, delete. That
// means five handlers, five route registrations, and five sets of
// JSON serialisation logic -- per model. Auto CRUD eliminates all
// of that repetition. You implement a ResourceService interface
// once, call kruda.Resource(), and the framework generates
// type-safe CRUD endpoints in a single call.

// Product represents a product in our catalogue.
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

// ============================================================
// ProductService -- ResourceService Implementation
// ============================================================
//
// kruda.ResourceService[T, ID] is the interface that Auto CRUD
// programs against. You implement 5 methods:
//
//   List(ctx, page, limit) ([]T, total, error)
//   Create(ctx, item)      (T, error)
//   Get(ctx, id)           (T, error)
//   Update(ctx, id, item)  (T, error)
//   Delete(ctx, id)        error
//
// The framework calls these methods from the generated handlers.
// You can put validation, business rules, or database access
// inside each method -- Auto CRUD doesn't care how you store
// the data, only that you satisfy the interface.

// ProductService is an in-memory implementation of
// ResourceService[Product, int]. In production you would
// replace this with a database-backed implementation without
// changing any route registration code.
type ProductService struct {
	mu     sync.Mutex
	items  []Product
	nextID int
}

// NewProductService creates a ProductService with an empty store.
func NewProductService() *ProductService {
	return &ProductService{nextID: 1}
}

// List returns a paginated slice of products and the total count.
// The page and limit parameters are extracted from query strings
// by the framework automatically (?page=1&limit=20).
func (s *ProductService) List(_ context.Context, page, limit int) ([]Product, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	total := len(s.items)

	// Calculate pagination bounds.
	start := (page - 1) * limit
	if start >= total {
		return nil, total, nil
	}
	end := start + limit
	if end > total {
		end = total
	}

	// Return a copy of the slice to prevent mutation.
	result := make([]Product, end-start)
	copy(result, s.items[start:end])
	return result, total, nil
}

// Create validates and persists a new product.
// Validation logic lives here -- if the input is invalid,
// return an error and the framework responds with 400/422.
func (s *ProductService) Create(_ context.Context, item Product) (Product, error) {
	if item.Name == "" {
		return Product{}, fmt.Errorf("product name is required")
	}
	if item.Price <= 0 {
		return Product{}, fmt.Errorf("price must be greater than zero")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	item.ID = s.nextID
	s.nextID++
	s.items = append(s.items, item)
	return item, nil
}

// Get returns a single product by ID.
func (s *ProductService) Get(_ context.Context, id int) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range s.items {
		if p.ID == id {
			return p, nil
		}
	}
	return Product{}, fmt.Errorf("product with id %d not found", id)
}

// Update replaces a product's data by ID.
func (s *ProductService) Update(_ context.Context, id int, item Product) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.items {
		if p.ID == id {
			item.ID = id // Preserve the original ID.
			s.items[i] = item
			return item, nil
		}
	}
	return Product{}, fmt.Errorf("product with id %d not found", id)
}

// Delete removes a product by ID.
func (s *ProductService) Delete(_ context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.items {
		if p.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("product with id %d not found", id)
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	app := kruda.New()

	// Create the product service.
	svc := NewProductService()

	// Register Auto CRUD for the Product model.
	//
	// kruda.Resource is the core of this tutorial. A single call
	// generates five fully functional, type-safe endpoints:
	//
	//   GET    /products      -- list all products (with pagination)
	//   POST   /products      -- create a new product
	//   GET    /products/:id  -- get a product by ID
	//   PUT    /products/:id  -- update a product by ID
	//   DELETE /products/:id  -- delete a product by ID
	//
	// Compare this to the beginner tutorial where we wrote each
	// handler by hand. Resource gives you the same result with
	// a fraction of the code.
	//
	// Optional ResourceOptions:
	//   kruda.WithResourceMiddleware(mw...)  -- add middleware
	//   kruda.WithResourceOnly("GET","POST") -- register only specific methods
	//   kruda.WithResourceExcept("DELETE")   -- exclude specific methods
	//   kruda.WithResourceIDParam("product_id") -- custom ID param name
	kruda.Resource[Product, int](app, "/products", svc)

	// Start the server on port 3000.
	log.Println("Server starting on :3000 ...")
	log.Fatal(app.Listen(":3000"))
}
