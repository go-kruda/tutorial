# 📄 Section 04-03 — OpenAPI Generator: Automatic API Documentation

⏱️ Estimated time: **30 minutes**

Welcome to the OpenAPI Generator lesson! 🎉 In this lesson you'll learn how to **automatically generate API documentation** from Typed Handlers using Kruda's **OpenAPI Generator** — no more writing YAML or JSON schemas by hand, because Kruda reads type information from Go generics and generates an OpenAPI 3.0 spec for you automatically.

---

## 🎯 Lesson Objectives

- Understand the concept of **OpenAPI Generator** and why it matters
- Learn how to enable OpenAPI with `kruda.WithOpenAPIInfo()` and `kruda.WithOpenAPITag()`
- Assign metadata to routes with `kruda.WithTags()` and `kruda.WithDescription()`
- See how Typed Handlers are converted into an OpenAPI spec
- Access API documentation via `/docs` and `/openapi.json`

---

## 📚 What You'll Learn

By the end of this lesson you'll be able to:

- ✅ Enable OpenAPI with `kruda.WithOpenAPIInfo("title", "version", "description")`
- ✅ Set the API title, description, and version
- ✅ Use `kruda.WithOpenAPITag()` to define tag group descriptions
- ✅ Use `kruda.WithTags()` to group endpoints in the spec
- ✅ Use `kruda.WithDescription()` to describe each operation
- ✅ Access the OpenAPI spec at `/openapi.json`
- ✅ View interactive API docs at `/docs` (Swagger UI)
- ✅ Understand how request/response types are automatically converted to JSON Schema

---

## 📋 Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25 or higher |
| Git | Latest version |
| Text Editor / IDE | VS Code, GoLand, or your preferred editor |
| Terminal | For running `go run` |
| Browser | For viewing Swagger UI at `/docs` |

> 💡 If you haven't completed Section 04-02 (Auth Middleware) yet, it's recommended to go back and do it first — 👉 [Section 04-02 — Auth Middleware](../02-auth-middleware/)

---

## 📂 File Structure

```
04-advanced/03-openapi/
├── README.md          ← 📖 You are here
├── starter/           ← 🏗️ Starter code (with TODOs to fill in)
│   ├── go.mod
│   └── main.go
└── complete/          ← ✅ Fully working reference implementation
    ├── go.mod
    └── main.go
```

- 📁 **[starter/](./starter/)** — A compilable skeleton with `// TODO:` placeholders for you to fill in
- 📁 **[complete/](./complete/)** — A fully working reference implementation you can compare against

---

## 💡 Why Use OpenAPI Generator?

Imagine you have an API with 30 endpoints — if you write the OpenAPI spec by hand, you'll run into problems:

| Problem | How OpenAPI Generator Helps |
|---|---|
| Spec doesn't match actual code (drift) | Generated directly from Typed Handlers — drift is impossible |
| Redundant JSON Schema writing | Reads Go struct tags and generates schemas automatically |
| Forgetting to document new endpoints | Every registered route appears in the spec immediately |
| Having to maintain a separate spec file | No spec file needed — Kruda generates it at runtime |

> 📄 Kruda's OpenAPI Generator uses information from **Typed Handlers** (Go generics) to create request/response schemas and uses route metadata to generate the operations list

```
Typed Handlers + Route Metadata → OpenAPI 3.0 Spec → /openapi.json + /docs
```

---

## 🛠️ Step-by-Step Guide

### Step 1: Open the starter project

```bash
cd 04-advanced/03-openapi/starter
```

Open the `main.go` file — you'll see a structure with `// TODO:` comments indicating where you need to add code.

### Step 2: Understand the structure

In this lesson we'll build an API for managing Products and Orders with OpenAPI documentation:

```
Product Endpoints:  GET /products, POST /products, GET /products/:id
Order Endpoints:    GET /orders, POST /orders, GET /orders/:id
Documentation:      GET /openapi.json, GET /docs
```

- **Request/Response Types** — Go structs that Kruda uses to generate JSON Schema
- **Typed Handlers** — Handlers with type information for OpenAPI
- **Route Metadata** — Tags and descriptions that appear in the spec

### Step 3: Create the Kruda App with OpenAPI Config

Replace the `// TODO:` for OpenAPI configuration in `main()`:

```go
app := kruda.New(
    kruda.WithOpenAPIInfo(
        "Product & Order API",
        "1.0.0",
        "A sample API demonstrating Kruda's OpenAPI Generator with Typed Handlers",
    ),
    kruda.WithOpenAPITag("Products", "Product management operations"),
    kruda.WithOpenAPITag("Orders", "Order management operations"),
)
```

> 📄 `WithOpenAPIInfo()` sets the metadata (title, version, description) that appears at the top of the OpenAPI 3.0 spec — while `WithOpenAPITag()` defines the description for each tag group in the spec

### Step 4: Register Product Routes with Metadata

Replace the `// TODO:` for Product route registration — use `kruda.Get[]` and `kruda.Post[]` with inline handlers:

```go
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
```

> 🧩 Notice that `kruda.Get[struct{}, []ProductResponse]` tells Kruda the input is `struct{}` (none) and the output is `[]ProductResponse` — this information will appear in the OpenAPI spec automatically

> 🏷️ `WithTags("Products")` groups endpoints in Swagger UI — making the API documentation easier to read

> 📝 `WithDescription("...")` adds a short description to each operation in the spec

### Step 5: Register Order Routes

Do the same for Order routes:

```go
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
```

> 🔍 Notice that the `POST /orders` handler uses `kruda.NotFound()` to return an error when the product is not found — Kruda will automatically convert it to an HTTP 404 response

### Step 6: Run and test

```bash
go run main.go
```

Open your browser and go to:

- 📄 **http://localhost:3000/openapi.json** — View the raw OpenAPI spec (JSON)
- 🖥️ **http://localhost:3000/docs** — View the interactive Swagger UI

Test the API with `curl`:

```bash
# Create a product
curl -X POST http://localhost:3000/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Kruda T-Shirt","description":"Official merch","price":590,"category":"apparel"}'

# List all products
curl http://localhost:3000/products

# Create an order
curl -X POST http://localhost:3000/orders \
  -H "Content-Type: application/json" \
  -d '{"product_id":1,"quantity":2}'

# View the OpenAPI spec
curl http://localhost:3000/openapi.json
```

If everything works correctly you'll see:
- Product and Order endpoints working normally
- The OpenAPI spec at `/openapi.json` with complete schemas for every type
- Swagger UI at `/docs` showing endpoints grouped by tags 🎉

---

## 🔍 Compare with complete/

If you get stuck anywhere, check the reference implementation in **[complete/](./complete/)** and compare it with your code, or use the `diff` command:

```bash
diff starter/main.go complete/main.go
```

---

## 💡 Key Takeaways

| Concept | Description |
|---|---|
| `kruda.WithOpenAPIInfo()` | Option in `kruda.New()` for setting title, version, description |
| `kruda.WithOpenAPITag()` | Option in `kruda.New()` for defining tag group descriptions |
| `kruda.WithTags()` | Groups endpoints in the OpenAPI spec |
| `kruda.WithDescription()` | Adds a short description to each operation |
| `kruda.Get[In, Out]()` | Registers a GET route with an inline typed handler |
| `kruda.Post[In, Out]()` | Registers a POST route with an inline typed handler |
| `*kruda.C[T]` | Context with typed input (`c.In`) for accessing request data |
| `kruda.NotFound()` | Creates an error that Kruda converts to an HTTP 404 response |
| `/openapi.json` | Endpoint for viewing the raw OpenAPI 3.0 spec |
| `/docs` | Endpoint for viewing the interactive Swagger UI |

---

## ➡️ Next Lesson

Awesome! 🎊 You've learned how to automatically generate API documentation with Kruda's OpenAPI Generator. In the next lesson we'll learn about **SSE (Server-Sent Events)** — how to push real-time data from server to client 📡

👉 [Section 04-04 — SSE](../04-sse/)
