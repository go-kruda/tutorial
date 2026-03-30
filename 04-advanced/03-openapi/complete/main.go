package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Request / Response Types
// ============================================================
//
// Kruda's OpenAPI Generator inspects these struct types and
// their tags to produce an accurate OpenAPI 3.0 spec automatically.
// The validate:"..." tags become JSON Schema constraints in the
// generated spec (e.g., required, minLength, minimum).

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

// ============================================================
// In-Memory Stores
// ============================================================

var (
	products    []ProductResponse
	orders      []OrderResponse
	mu          sync.Mutex
	nextProdID  = 1
	nextOrderID = 1
)

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	// ── 1. Create the Kruda Application with OpenAPI config ──
	//
	// WithOpenAPIInfo configures the metadata that appears at the
	// top of the generated OpenAPI 3.0 specification.
	// WithOpenAPITag defines tag descriptions for grouping.
	app := kruda.New(
		kruda.WithOpenAPIInfo(
			"Product & Order API",
			"1.0.0",
			"A sample API demonstrating Kruda's OpenAPI Generator with Typed Handlers",
		),
		kruda.WithOpenAPITag("Products", "Product management operations"),
		kruda.WithOpenAPITag("Orders", "Order management operations"),
	)

	// ── 2. Register Product Routes ───────────────────────────
	//
	// Each route registration feeds into the OpenAPI spec.
	// Kruda inspects the typed handler's generic parameters to
	// determine request/response schemas. WithDescription() and
	// WithTags() add metadata to each operation.

	kruda.Get[struct{}, []ProductResponse](app, "/products",
		func(c *kruda.C[struct{}]) (*[]ProductResponse, error) {
			mu.Lock()
			defer mu.Unlock()
			result := make([]ProductResponse, len(products))
			copy(result, products)
			return &result, nil
		},
		kruda.WithDescription("List all products"),
		kruda.WithTags("Products"),
	)

	kruda.Post[CreateProductInput, ProductResponse](app, "/products",
		func(c *kruda.C[CreateProductInput]) (*ProductResponse, error) {
			mu.Lock()
			defer mu.Unlock()
			product := ProductResponse{
				ID:          nextProdID,
				Name:        c.In.Name,
				Description: c.In.Description,
				Price:       c.In.Price,
				Category:    c.In.Category,
			}
			nextProdID++
			products = append(products, product)
			return &product, nil
		},
		kruda.WithDescription("Create a new product"),
		kruda.WithTags("Products"),
	)

	kruda.Get[GetProductInput, ProductResponse](app, "/products/:id",
		func(c *kruda.C[GetProductInput]) (*ProductResponse, error) {
			mu.Lock()
			defer mu.Unlock()
			for _, p := range products {
				if p.ID == c.In.ID {
					return &p, nil
				}
			}
			return nil, kruda.NotFound(fmt.Sprintf("product %d not found", c.In.ID))
		},
		kruda.WithDescription("Get a product by ID"),
		kruda.WithTags("Products"),
	)

	// ── 3. Register Order Routes ─────────────────────────────

	kruda.Get[struct{}, []OrderResponse](app, "/orders",
		func(c *kruda.C[struct{}]) (*[]OrderResponse, error) {
			mu.Lock()
			defer mu.Unlock()
			result := make([]OrderResponse, len(orders))
			copy(result, orders)
			return &result, nil
		},
		kruda.WithDescription("List all orders"),
		kruda.WithTags("Orders"),
	)

	kruda.Post[CreateOrderInput, OrderResponse](app, "/orders",
		func(c *kruda.C[CreateOrderInput]) (*OrderResponse, error) {
			mu.Lock()
			defer mu.Unlock()

			var price float64
			found := false
			for _, p := range products {
				if p.ID == c.In.ProductID {
					price = p.Price
					found = true
					break
				}
			}
			if !found {
				return nil, kruda.NotFound(fmt.Sprintf("product %d not found", c.In.ProductID))
			}

			order := OrderResponse{
				ID:        nextOrderID,
				ProductID: c.In.ProductID,
				Quantity:  c.In.Quantity,
				Total:     price * float64(c.In.Quantity),
				Status:    "pending",
			}
			nextOrderID++
			orders = append(orders, order)
			return &order, nil
		},
		kruda.WithDescription("Create a new order"),
		kruda.WithTags("Orders"),
	)

	kruda.Get[GetOrderInput, OrderResponse](app, "/orders/:id",
		func(c *kruda.C[GetOrderInput]) (*OrderResponse, error) {
			mu.Lock()
			defer mu.Unlock()
			for _, o := range orders {
				if o.ID == c.In.ID {
					return &o, nil
				}
			}
			return nil, kruda.NotFound(fmt.Sprintf("order %d not found", c.In.ID))
		},
		kruda.WithDescription("Get an order by ID"),
		kruda.WithTags("Orders"),
	)

	// ── 4. Start the Server ──────────────────────────────────
	//
	// Kruda automatically serves the OpenAPI spec at:
	//   GET /openapi.json  -- the raw JSON spec
	//   GET /docs          -- Swagger UI viewer
	log.Println("Server starting on :3000 ...")
	log.Println("  OpenAPI: GET /openapi.json")
	log.Println("  Docs:    GET /docs")
	log.Fatal(app.Listen(":3000"))
}
