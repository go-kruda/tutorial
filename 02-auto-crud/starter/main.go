package main

import (
	"context"
	"log"
	"sync"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Model Definition
// ============================================================

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
// TODO: implement kruda.ResourceService[Product, int] interface
//
// Interface to implement:
//   List(ctx context.Context, page, limit int) ([]Product, int, error)
//   Create(ctx context.Context, item Product) (Product, error)
//   Get(ctx context.Context, id int) (Product, error)
//   Update(ctx context.Context, id int, item Product) (Product, error)
//   Delete(ctx context.Context, id int) error

type ProductService struct {
	mu     sync.Mutex
	items  []Product
	nextID int
}

func NewProductService() *ProductService {
	return &ProductService{nextID: 1}
}

// TODO: implement List -- fetch product list with pagination
//
// Hint: Calculate start/end from page and limit
//   start := (page - 1) * limit
//   return a slice of items[start:end] and total count
func (s *ProductService) List(_ context.Context, page, limit int) ([]Product, int, error) {
	// TODO: implement pagination logic
	return nil, 0, nil
}

// TODO: implement Create -- create a new product with validation
//
// Hint: Validate that item.Name != "" and item.Price > 0
//   Assign ID automatically from s.nextID
func (s *ProductService) Create(_ context.Context, item Product) (Product, error) {
	// TODO: validate and persist
	return Product{}, nil
}

// TODO: implement Get -- fetch a product by ID
//
// Hint: Loop to find the item with a matching ID
//   If not found, return error
func (s *ProductService) Get(_ context.Context, id int) (Product, error) {
	// TODO: find by ID
	return Product{}, nil
}

// TODO: implement Update -- update a product by ID
//
// Hint: Find the item with a matching ID and replace it with the new values
//   Don't forget item.ID = id to preserve the original ID
func (s *ProductService) Update(_ context.Context, id int, item Product) (Product, error) {
	// TODO: find and update
	return Product{}, nil
}

// TODO: implement Delete -- delete a product by ID
//
// Hint: Find the item with a matching ID and remove it from the slice
func (s *ProductService) Delete(_ context.Context, id int) error {
	// TODO: find and delete
	return nil
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	app := kruda.New()
	svc := NewProductService()

	// TODO: Register Auto CRUD for Product
	//
	// Hint: Use kruda.Resource[T, ID] to automatically generate 5 endpoints
	//
	//   kruda.Resource[Product, int](app, "/products", svc)
	//
	// A single line gives you:
	//   GET    /products      -- list (with pagination)
	//   POST   /products      -- create
	//   GET    /products/:id  -- get by ID
	//   PUT    /products/:id  -- update
	//   DELETE /products/:id  -- delete
	_ = svc

	log.Println("Server starting on :3000 ...")
	log.Fatal(app.Listen(":3000"))
}
