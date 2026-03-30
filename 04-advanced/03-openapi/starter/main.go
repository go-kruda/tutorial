package main

import (
	"log"
	"sync"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Request / Response Types
// ============================================================

type CreateProductInput struct {
	Name        string  `json:"name" validate:"required,min=1,max=200"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Category    string  `json:"category"`
}

type GetProductInput struct {
	ID int `param:"id"`
}

type ProductResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

type CreateOrderInput struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,gt=0"`
}

type GetOrderInput struct {
	ID int `param:"id"`
}

type OrderResponse struct {
	ID        int     `json:"id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
	Status    string  `json:"status"`
}

var (
	products    []ProductResponse
	orders      []OrderResponse
	mu          sync.Mutex
	nextProdID  = 1
	nextOrderID = 1
)

func main() {
	// TODO: Create a Kruda app with OpenAPI config
	//
	// Hint: Use options in kruda.New():
	//   app := kruda.New(
	//       kruda.WithOpenAPIInfo("title", "1.0.0", "description"),
	//       kruda.WithOpenAPITag("Products", "Product operations"),
	//       kruda.WithOpenAPITag("Orders", "Order operations"),
	//   )
	app := kruda.New()

	// TODO: Register Product routes with OpenAPI metadata
	//
	// Hint: Use kruda.WithDescription() and kruda.WithTags() as route options
	//
	//   kruda.Get[struct{}, []ProductResponse](app, "/products",
	//       handler,
	//       kruda.WithDescription("List all products"),
	//       kruda.WithTags("Products"),
	//   )

	// TODO: Register Order routes with OpenAPI metadata

	_ = &mu
	_ = products
	_ = orders
	_ = nextProdID
	_ = nextOrderID

	log.Println("Server starting on :3000 ...")
	log.Fatal(app.Listen(":3000"))
}
