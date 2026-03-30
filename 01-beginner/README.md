# 🚀 Section 01 — Getting Started with Kruda: Basic REST API + Typed Handlers

⏱️ Estimated time: **30 minutes**

Welcome to the first Kruda lesson! 🎉 In this section you will build a REST API for managing books (Book API) from scratch using **Typed Handlers** — a standout Kruda feature that lets you write type-safe handlers with Go generics — no more manual JSON encode/decode!

---

## 🎯 Lesson Objectives

- Create your first Kruda application
- Understand the concept of **Typed Handler** with `kruda.Get[In, Out]()` / `kruda.Post[In, Out]()` / `kruda.Delete[In, Out]()`
- Build a REST API with 4 endpoints: `GET`, `POST`, `GET /:id`, `DELETE /:id`
- Learn how to use path parameters via the `param:"id"` struct tag (auto-parse)
- See how Kruda handles JSON serialisation/deserialisation automatically

---

## 📚 What You Will Learn

By the end of this lesson you will be able to:

- ✅ Create a Kruda app with `kruda.New()`
- ✅ Define input/output types as Go structs
- ✅ Write Typed Handlers with `kruda.Get[In, Out]()`, `kruda.Post[In, Out]()`, `kruda.Delete[In, Out]()`
- ✅ Automatically extract path parameters with the `param:"id"` struct tag and `c.In.ID`
- ✅ Use the `*kruda.C[T]` context which contains parsed input in `c.In`
- ✅ Run the server with `app.Listen(":3000")`

---

## 📋 Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25 or higher |
| Git | Latest version |
| Text Editor / IDE | VS Code, GoLand, or your preferred editor |
| Terminal | For running `go run` |

> 💡 If you have never written Go before, we recommend trying [A Tour of Go](https://go.dev/tour/) first

---

## 📂 File Structure

```
01-beginner/
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
cd 01-beginner/starter
```

Open the `main.go` file — you will see a structure with `// TODO:` comments indicating where you need to add code.

### Step 2: Define Input/Output Types

The starter file already has these structs prepared for you:

```go
type CreateBookInput struct {
    Title  string `json:"title"`
    Author string `json:"author"`
}

type GetBookByIDInput struct {
    ID int `param:"id"`
}

type BookResponse struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

type MessageResponse struct {
    Message string `json:"message"`
}
```

> 🧩 Kruda uses struct tags (`json:"..."`) to automatically map between Go fields and JSON keys, and `param:"id"` to auto-parse path parameters from the URL — no need to call `strconv.Atoi` yourself!

### Step 3: Create an In-Memory Store

Add variables to store book data in memory:

```go
var (
    books  []BookResponse
    mu     sync.Mutex
    nextID = 1
)
```

We use `sync.Mutex` to prevent race conditions when multiple requests come in simultaneously.

### Step 4: Write Your First Typed Handler — GET /books

This is the heart of the lesson! Register a handler for `GET /books`:

```go
kruda.Get[struct{}, []BookResponse](app, "/books", func(c *kruda.C[struct{}]) (*[]BookResponse, error) {
    mu.Lock()
    defer mu.Unlock()

    result := make([]BookResponse, len(books))
    copy(result, books)
    return &result, nil
})
```

> 🎯 Notice that the input type is `struct{}` because `GET /books` has no request body or path parameter — Kruda will automatically serialise `[]BookResponse` into a JSON array. The handler returns a pointer `*[]BookResponse`

### Step 5: Write the Handler for POST /books

```go
kruda.Post[CreateBookInput, BookResponse](app, "/books", func(c *kruda.C[CreateBookInput]) (*BookResponse, error) {
    mu.Lock()
    defer mu.Unlock()

    book := BookResponse{
        ID:     nextID,
        Title:  c.In.Title,
        Author: c.In.Author,
    }
    nextID++
    books = append(books, book)
    return &book, nil
})
```

> ✨ Kruda automatically deserialises the JSON body into `CreateBookInput` — access values via `c.In.Title` and `c.In.Author` without calling `json.NewDecoder` yourself!

### Step 6: Write the Handler for GET /books/:id

```go
kruda.Get[GetBookByIDInput, BookResponse](app, "/books/:id", func(c *kruda.C[GetBookByIDInput]) (*BookResponse, error) {
    mu.Lock()
    defer mu.Unlock()

    for _, b := range books {
        if b.ID == c.In.ID {
            return &b, nil
        }
    }

    return nil, kruda.NotFound(fmt.Sprintf("book with id %d not found", c.In.ID))
})
```

> 🔑 `GetBookByIDInput` has a field `ID int` with the tag `param:"id"` — Kruda will parse `:id` from the URL and convert it to `int` automatically. Access the value via `c.In.ID` without calling `strconv.Atoi` yourself! Use `kruda.NotFound()` for HTTP error responses

### Step 7: Write the Handler for DELETE /books/:id

```go
kruda.Delete[GetBookByIDInput, MessageResponse](app, "/books/:id", func(c *kruda.C[GetBookByIDInput]) (*MessageResponse, error) {
    mu.Lock()
    defer mu.Unlock()

    for i, b := range books {
        if b.ID == c.In.ID {
            books = append(books[:i], books[i+1:]...)
            return &MessageResponse{
                Message: fmt.Sprintf("book %d deleted", c.In.ID),
            }, nil
        }
    }

    return nil, kruda.NotFound(fmt.Sprintf("book with id %d not found", c.In.ID))
})
```

### Step 8: Create the App and Register Routes

Put everything together in `main()`:

```go
func main() {
    app := kruda.New()

    // Add middleware: Recovery catches panics, Logger logs every request
    app.Use(middleware.Recovery(), middleware.Logger())

    // Register routes with kruda.Get, kruda.Post, kruda.Delete
    // (See Steps 4-7 for each handler)

    kruda.Get[struct{}, []BookResponse](app, "/books", func(c *kruda.C[struct{}]) (*[]BookResponse, error) {
        // ... handler code ...
    })

    kruda.Post[CreateBookInput, BookResponse](app, "/books", func(c *kruda.C[CreateBookInput]) (*BookResponse, error) {
        // ... handler code ...
    })

    kruda.Get[GetBookByIDInput, BookResponse](app, "/books/:id", func(c *kruda.C[GetBookByIDInput]) (*BookResponse, error) {
        // ... handler code ...
    })

    kruda.Delete[GetBookByIDInput, MessageResponse](app, "/books/:id", func(c *kruda.C[GetBookByIDInput]) (*MessageResponse, error) {
        // ... handler code ...
    })

    log.Println("Server starting on :3000 ...")
    log.Fatal(app.Listen(":3000"))
}
```

> 🚀 `kruda.New()` creates an app with Wing Transport powered by epoll — you get high-performance networking without any extra configuration. Notice that route registration uses package-level functions `kruda.Get()`, `kruda.Post()`, `kruda.Delete()` instead of methods on the app

### Step 9: Run and Test

```bash
go run main.go
```

Open another terminal window and test with `curl`:

```bash
# Create the first book
curl -X POST http://localhost:3000/books \
  -H "Content-Type: application/json" \
  -d '{"title":"Go Programming","author":"John Doe"}'

# List all books
curl http://localhost:3000/books

# Get a book by ID
curl http://localhost:3000/books/1

# Delete a book
curl -X DELETE http://localhost:3000/books/1
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
| `kruda.New()` | Creates a Kruda app instance with Wing Transport |
| `middleware.Recovery()` | Middleware that catches panics and returns 500 instead of crashing |
| `middleware.Logger()` | Middleware that logs every request (method, path, status, latency) |
| `app.Use(mw...)` | Adds global middleware to all routes |
| `kruda.Get[In, Out]()` / `kruda.Post[In, Out]()` / `kruda.Delete[In, Out]()` | Registers a Typed Handler with compile-time input/output types |
| `*kruda.C[T]` | Context with parsed input in `c.In` |
| `param:"id"` struct tag | Tells Kruda to auto-parse the path parameter from the URL into the specified type |
| `c.In` | Access parsed input (JSON body, path params, query params) |
| `kruda.NotFound()` | Creates an HTTP error response for 404 Not Found |
| `app.Listen(":3000")` | Starts the server on the specified port |

---

## ➡️ Next Section

Ready? Let's move on! In the next section we will learn about **Auto CRUD** — a feature that automatically generates CRUD endpoints from a model struct, massively reducing repetitive code 🔥

👉 [Section 02 — Auto CRUD](../02-auto-crud/)
