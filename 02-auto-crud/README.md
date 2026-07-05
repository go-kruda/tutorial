# 🤖 Section 02 — Auto CRUD: Automatically Generate CRUD Endpoints with ResourceService

⏱️ Estimated time: **30 minutes**

Welcome to the second lesson! 🎉 In this section you will learn about **Auto CRUD** — a standout Kruda feature that generates all 5 CRUD endpoints (List, Create, Get, Update, Delete) from a single `ResourceService` interface, eliminating the need to write repetitive handlers ever again.

---

## 🎯 Lesson Objectives

- Understand what **Auto CRUD** is and why it massively reduces boilerplate code
- Learn how to define a model struct for Auto CRUD
- Create a service that implements the `ResourceService[T, ID]` interface
- Register Auto CRUD with a single call to `kruda.Resource[T, ID]()`
- Understand optional `ResourceOptions` for middleware, method filtering, and custom ID param

---

## 📚 What You Will Learn

By the end of this lesson you will be able to:

- ✅ Explain how Auto CRUD works and how it differs from writing handlers manually
- ✅ Define a model struct for Auto CRUD
- ✅ Implement the `ResourceService[T, ID]` interface with 5 methods: `List`, `Create`, `Get`, `Update`, `Delete`
- ✅ Add automatic validation with `validate` tags + `kruda.WithValidator(...)`, and business logic in service methods
- ✅ Return `kruda.NotFound()` from a service method to get a real 404
- ✅ Register a model with `kruda.Resource[T, ID]()` to automatically generate 5 endpoints

---

## 💡 What is Auto CRUD?

In the previous lesson (Section 01) you wrote a separate handler for each endpoint — `GET /books`, `POST /books`, `GET /books/:id`, `DELETE /books/:id` — totalling 4 handlers plus 4 lines of route registration. Imagine if your app had 10 models — you would have to write the same repetitive code 10 times 😵

**Auto CRUD** solves this problem by letting you:

1. Define a model struct (e.g. `Product`)
2. Create a service that implements the `ResourceService[T, ID]` interface
3. Call `kruda.Resource[T, ID]()` just once

Kruda will automatically generate all 5 endpoints:

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/products` | List all items (with pagination support) |
| `POST` | `/products` | Create a new item |
| `GET` | `/products/:id` | Get an item by ID |
| `PUT` | `/products/:id` | Update an item by ID |
| `DELETE` | `/products/:id` | Delete an item by ID |

> 🔥 Compare this with Section 01 where you had to write 4 handlers manually — Auto CRUD delivers all 5 endpoints with just a few lines of code!

---

## 📋 Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25 or higher |
| Git | Latest version |
| Text Editor / IDE | VS Code, GoLand, or your preferred editor |
| Terminal | For running `go run` |

> 💡 If you haven't completed Section 01 yet, we recommend going back to do it first since this section builds on the concepts of Typed Handlers — 👉 [Section 01 — Beginner](../01-beginner/)

---

## 📂 File Structure

```
02-auto-crud/
├── README.md          ← 📖 You are here
├── starter/           ← 🏗️ Starter code (with TODOs to fill in)
│   ├── go.mod
│   └── main.go
└── complete/          ← ✅ Fully working reference implementation
    ├── go.mod
    └── main.go
```

- 📁 **[starter/](./starter/)** — A skeleton that compiles but still has `// TODO:` placeholders for you to fill in
- 📁 **[complete/](./complete/)** — A fully working reference implementation you can compare against

---

## 🛠️ Step-by-Step

### Step 1: Open the starter project

```bash
cd 02-auto-crud/starter
```

Open the `main.go` file — you will see a structure with `// TODO:` comments indicating where you need to add code.

### Step 2: Understand the Model Struct

The starter file already has a `Product` struct prepared for you:

```go
type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name" validate:"required"`
    Price float64 `json:"price" validate:"required,gt=0"`
    Stock int     `json:"stock"`
}
```

> 🧩 The model struct uses `json` tags for JSON serialisation and `validate` tags for automatic input validation — Kruda will use this struct as a type parameter in `kruda.Resource[Product, int]()` to generate type-safe, validated endpoints

### Step 3: Understand the ResourceService Interface

Auto CRUD works through the `ResourceService[T, ID]` interface which requires you to implement 5 methods:

```go
// kruda.ResourceService[T, ID] interface:
//   List(ctx context.Context, page, limit int) ([]T, int, error)
//   Create(ctx context.Context, item T) (T, error)
//   Get(ctx context.Context, id ID) (T, error)
//   Update(ctx context.Context, id ID, item T) (T, error)
//   Delete(ctx context.Context, id ID) error
```

> 💡 The framework will call these methods from the generated handlers — you can put validation, business rules, or database access in the methods. Auto CRUD doesn't care how you store data, just implement the interface completely

### Step 4: Create the ProductService

Create an in-memory implementation of `ResourceService[Product, int]`:

```go
type ProductService struct {
    mu     sync.Mutex
    items  []Product
    nextID int
}

func NewProductService() *ProductService {
    return &ProductService{nextID: 1}
}
```

### Step 5: Implement the List Method

```go
func (s *ProductService) List(_ context.Context, page, limit int) ([]Product, int, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    total := len(s.items)

    start := (page - 1) * limit
    if start >= total {
        return nil, total, nil
    }
    end := start + limit
    if end > total {
        end = total
    }

    result := make([]Product, end-start)
    copy(result, s.items[start:end])
    return result, total, nil
}
```

> 🎯 `page` and `limit` are automatically extracted from query strings (`?page=1&limit=20`) by the framework — you just receive and use the values

### Step 6: Implement the Create Method

Since kruda **v1.4.0**, `kruda.Resource` validates create/update request bodies automatically -- using the `validate` tags on your model (Step 2) plus `kruda.WithValidator(...)` (Step 8) -- and returns a `422` *before* your service method is even called. `Create` only needs to persist the item:

```go
func (s *ProductService) Create(_ context.Context, item Product) (Product, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    item.ID = s.nextID
    s.nextID++
    s.items = append(s.items, item)
    return item, nil
}
```

> 🔧 Field-shape checks (required, format, ranges) belong on the model as `validate` tags -- the framework enforces them for you. Keep manual checks in the service only for cross-field or business rules that can't be expressed as a tag (there aren't any in this example). Without `kruda.WithValidator(...)` set on the app, `validate` tags do nothing -- Step 8 wires it up.

### Step 7: Implement Get, Update, Delete Methods

```go
func (s *ProductService) Get(_ context.Context, id int) (Product, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    for _, p := range s.items {
        if p.ID == id {
            return p, nil
        }
    }
    return Product{}, kruda.NotFound(fmt.Sprintf("product with id %d not found", id))
}

func (s *ProductService) Update(_ context.Context, id int, item Product) (Product, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    for i, p := range s.items {
        if p.ID == id {
            item.ID = id
            s.items[i] = item
            return item, nil
        }
    }
    return Product{}, kruda.NotFound(fmt.Sprintf("product with id %d not found", id))
}

func (s *ProductService) Delete(_ context.Context, id int) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    for i, p := range s.items {
        if p.ID == id {
            s.items = append(s.items[:i], s.items[i+1:]...)
            return nil
        }
    }
    return kruda.NotFound(fmt.Sprintf("product with id %d not found", id))
}
```

> ⚠️ Fix: these three methods used to return a plain `fmt.Errorf(...)` on a missing ID, which `kruda.Resource` resolves to a generic `500` -- not the `404` you'd expect. Returning `kruda.NotFound(...)` (a `*KrudaError`) flows through unchanged and renders as a real `404`.

### Step 8: Create the App and Register Auto CRUD

This is the heart of the lesson! Replace the `// TODO:` in `main()`:

```go
func main() {
    // kruda.WithValidator(kruda.NewValidator()) activates the `validate`
    // tags on Product -- without it, the tags do nothing.
    app := kruda.New(kruda.WithValidator(kruda.NewValidator()))

    svc := NewProductService()

    // kruda.Resource is the heart of Auto CRUD
    // A single call generates 5 endpoints:
    //   GET    /products      -- list (with pagination)
    //   POST   /products      -- create
    //   GET    /products/:id  -- get by ID
    //   PUT    /products/:id  -- update by ID
    //   DELETE /products/:id  -- delete by ID
    kruda.Resource[Product, int](app, "/products", svc)

    log.Println("Server starting on :3000 ...")
    log.Fatal(app.Listen(":3000"))
}
```

> 🚀 Just a single call to `kruda.Resource[Product, int]()` and Kruda automatically generates all 5 endpoints — compare this with Section 01 where you had to write each handler and register each route individually!

Optional `ResourceOptions` you can use:

| Option | Description |
|--------|-------------|
| `kruda.WithResourceMiddleware(mw...)` | Add middleware for this resource |
| `kruda.WithResourceOnly("GET","POST")` | Register only the specified methods |
| `kruda.WithResourceExcept("DELETE")` | Exclude the specified methods |
| `kruda.WithResourceIDParam("product_id")` | Change the ID parameter name |

### Step 9: Run and Test

```bash
go run main.go
```

Open another terminal window and test with `curl`:

```bash
# Create a new product
curl -X POST http://localhost:3000/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Mechanical Keyboard","price":2500,"stock":50}'

# List all products
curl http://localhost:3000/products

# Get a product by ID
curl http://localhost:3000/products/1

# Update a product
curl -X PUT http://localhost:3000/products/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Mechanical Keyboard RGB","price":2900,"stock":45}'

# Delete a product
curl -X DELETE http://localhost:3000/products/1

# Test validation -- 422, name is required and price must be > 0
curl -X POST http://localhost:3000/products \
  -H "Content-Type: application/json" \
  -d '{"name":"","price":0,"stock":10}'

# Test 404 -- now a real 404, not 500
curl http://localhost:3000/products/999
```

If everything works correctly, you will see JSON responses returned for each command 🎉

---

## 🔍 Compare with complete/

If you get stuck anywhere, check the reference implementation in **[complete/](./complete/)** and compare it with your code, or use the `diff` command:

```bash
diff starter/main.go complete/main.go
```

---

## 💡 Key Concepts Summary

| Concept | Description |
|---|---|
| Auto CRUD | A feature that automatically generates CRUD endpoints from a `ResourceService` interface |
| `ResourceService[T, ID]` | An interface requiring 5 methods: `List`, `Create`, `Get`, `Update`, `Delete` |
| `kruda.Resource[T, ID]()` | Registers a service to generate 5 CRUD endpoints in a single line |
| `kruda.WithValidator(kruda.NewValidator())` | Activates `validate` struct tags -- required for Resource's auto-validation |
| `validate` struct tags | Declarative field-shape validation on your model (e.g. `required`, `gt=0`) -- Resource enforces them before your service runs |
| `WithResourceMiddleware` | Adds middleware for a resource |
| `WithResourceOnly` / `WithResourceExcept` | Include/exclude HTTP methods to generate |
| `kruda.NotFound()` in a service method | Returns a real 404 from `Get`/`Update`/`Delete` -- a plain error resolves to 500 |

---

## ➡️ Next Section

Awesome! 🎊 You have learned how to use Auto CRUD to reduce boilerplate code. In the next section we will step into the intermediate level — connecting to a real database, managing config with environment variables, and writing structured error handling 💪

👉 [Section 03 — Intermediate](../03-intermediate/)
