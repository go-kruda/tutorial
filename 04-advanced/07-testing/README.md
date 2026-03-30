# 🧪 Section 04-07 — Testing: Unit Tests with TestClient

⏱️ Estimated time: **30 minutes**

Welcome to the Testing lesson! In this lesson you will learn how to write **unit tests** for a Kruda API using `kruda.NewTestClient(app)` together with Go's standard `testing` package -- no external libraries needed. Test routes without starting a real HTTP server.

---

## Learning Objectives

- Understand how to test a Kruda API without starting a server
- Learn how to use `kruda.NewTestClient(app)` to create a test client
- Write **table-driven tests**, which is a standard Go pattern
- Test both success cases and error cases (400, 404)
- Separate tests between the data layer (store) and the HTTP layer

---

## What You Will Learn

By the end of this lesson you will be able to:

- Use `kruda.NewTestClient(app)` to create a test client for testing
- Call `client.Get("/path")`, `client.Post("/path", body)`, etc. to test routes
- Check `resp.StatusCode()` and `resp.JSON(&v)` for assertions
- Write table-driven tests with `t.Run()` sub-tests
- Test error cases (empty input, not found)
- Run tests with `go test -v ./...`

---

## Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25+ |
| Git | Latest |
| Text Editor / IDE | VS Code, GoLand, or your preferred editor |
| Terminal | For running `go test` |

> If you haven't completed Section 04-06 (WebSocket), it's recommended to go back and do it first -- [Section 04-06 -- WebSocket](../06-websocket/)

---

## File Structure

```
04-advanced/07-testing/
├── README.md              <- You are here
├── starter/               <- Starter code (with TODOs to fill in)
│   ├── go.mod
│   └── main.go
└── complete/              <- Complete working solution
    ├── go.mod
    ├── main.go            <- API implementation + setupApp()
    └── handler_test.go    <- Unit tests with TestClient
```

- **[starter/](./starter/)** -- Skeleton code that compiles but has `// TODO:` markers for you to complete
- **[complete/](./complete/)** -- Full working solution with test file included

---

## Why Write Unit Tests?

Unit tests help you:

| Benefit | Description |
|---|---|
| Prevent regression | When you change code, tests will immediately tell you if something breaks |
| Serve as documentation | Tests show how the API is expected to behave |
| Fast feedback loop | Run `go test` instantly without starting a server |
| Test in isolation | Test each route independently |

### Testing Pattern in Kruda

```
┌──────────────────┐     ┌───────────────────────────┐
│  Test Function   │ ──> │  kruda.NewTestClient(app)  │
│  (handler_test)  │     │  (test client)             │
└──────────────────┘     └───────────────────────────┘
         |                         |
         v                         v
┌──────────────────┐     ┌───────────────────────────┐
│  client.Get()    │ ──> │  Full Kruda pipeline       │
│  client.Post()   │     │  (routing + handlers)      │
└──────────────────┘     └───────────────────────────┘
         |
         v
   resp.StatusCode()
   resp.JSON(&v)
   (assert in test)
```

> `kruda.NewTestClient(app)` creates a test client that simulates HTTP requests -- you can call GET, POST, PATCH, DELETE without starting a real server. It tests routing, input binding, and handler logic all at once.

---

## Step-by-Step Guide

### Step 1: Open the starter project

```bash
cd 04-advanced/07-testing/starter
```

Open the `main.go` file -- you will see the skeleton with `// TODO:` comments.

### Step 2: Implement TaskStore Methods

Start with the data layer first -- replace the `// TODO:` in `Create()`, `FindByID()`, and `ToggleDone()`:

```go
func (s *TaskStore) Create(title, description string) TaskResponse {
    s.mu.Lock()
    defer s.mu.Unlock()

    task := TaskResponse{
        ID:          s.nextID,
        Title:       title,
        Description: description,
        Done:        false,
    }
    s.nextID++
    s.tasks = append(s.tasks, task)
    return task
}
```

### Step 3: Implement Route Handlers

Replace the `// TODO:` in `setupAppWithStore()` to register all 4 routes. For example, POST /tasks:

```go
kruda.Post[CreateTaskInput, TaskResponse](app, "/tasks", func(c *kruda.C[CreateTaskInput]) (*TaskResponse, error) {
    if c.In.Title == "" {
        return nil, kruda.BadRequest("title is required")
    }
    task := store.Create(c.In.Title, c.In.Description)
    return &task, nil
})
```

> Notice that the handler performs validation (title must not be empty) -- we will write tests for both success and error cases.

### Step 4: Create the test file

Create a `handler_test.go` file in the same directory as `main.go`:

```bash
touch handler_test.go
```

> Go convention: test files must end with `_test.go` and be in the same package.

### Step 5: Write Your First Test -- TestCreateTask

```go
package main

import (
    "testing"
    "github.com/go-kruda/kruda"
)

func TestCreateTask_Success(t *testing.T) {
    app := setupApp()
    client := kruda.NewTestClient(app)

    resp, err := client.Post("/tasks", map[string]string{
        "title":       "Write unit tests",
        "description": "Learn testing",
    })
    if err != nil {
        t.Fatalf("expected no error, got: %v", err)
    }
    if resp.StatusCode() != 200 {
        t.Errorf("expected status 200, got %d", resp.StatusCode())
    }

    var task TaskResponse
    resp.JSON(&task)

    if task.ID != 1 {
        t.Errorf("expected ID=1, got ID=%d", task.ID)
    }
    if task.Title != "Write unit tests" {
        t.Errorf("expected Title=%q, got %q", "Write unit tests", task.Title)
    }
}
```

> `client.Post("/tasks", body)` sends a request through the test client -- you get a response back with `StatusCode()` and `JSON()` for assertions.

### Step 6: Write Table-Driven Tests

Table-driven tests are a standard Go pattern for testing multiple scenarios:

```go
func TestCreateTask_TableDriven(t *testing.T) {
    tests := []struct {
        name       string
        title      string
        desc       string
        wantStatus int
        wantTitle  string
    }{
        {
            name:       "valid task",
            title:      "Buy groceries",
            desc:       "Milk, eggs",
            wantStatus: 200,
            wantTitle:  "Buy groceries",
        },
        {
            name:       "empty title returns 400",
            title:      "",
            desc:       "Should fail",
            wantStatus: 400,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            app := setupApp()
            client := kruda.NewTestClient(app)

            resp, err := client.Post("/tasks", map[string]string{
                "title":       tt.title,
                "description": tt.desc,
            })
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }

            if resp.StatusCode() != tt.wantStatus {
                t.Errorf("status = %d, want %d", resp.StatusCode(), tt.wantStatus)
            }

            if tt.wantStatus == 200 {
                var task TaskResponse
                resp.JSON(&task)
                if task.Title != tt.wantTitle {
                    t.Errorf("Title = %q, want %q", task.Title, tt.wantTitle)
                }
            }
        })
    }
}
```

> Each sub-test creates a fresh app for isolation -- one test will not affect another.

### Step 7: Test Path Parameters and Error Codes

Path parameters are automatically bound via the struct tag `param:"id"` -- in tests you simply call the URL with the ID:

```go
func TestGetTaskByID(t *testing.T) {
    app := setupApp()
    client := kruda.NewTestClient(app)

    // Create a task first
    client.Post("/tasks", map[string]string{"title": "Existing task"})

    // Get it by ID
    resp, _ := client.Get("/tasks/1")
    if resp.StatusCode() != 200 {
        t.Errorf("expected 200, got %d", resp.StatusCode())
    }

    var task TaskResponse
    resp.JSON(&task)
    if task.Title != "Existing task" {
        t.Errorf("Title = %q, want %q", task.Title, "Existing task")
    }
}

func TestGetTaskByID_NotFound(t *testing.T) {
    app := setupApp()
    client := kruda.NewTestClient(app)

    resp, _ := client.Get("/tasks/999")
    if resp.StatusCode() != 404 {
        t.Errorf("expected 404, got %d", resp.StatusCode())
    }
}
```

### Step 8: Run Tests

```bash
go test -v ./...
```

You will see output like this:

```
=== RUN   TestCreateTask_Success
--- PASS: TestCreateTask_Success (0.00s)
=== RUN   TestCreateTask_TableDriven
=== RUN   TestCreateTask_TableDriven/valid_task
=== RUN   TestCreateTask_TableDriven/empty_title_returns_400
--- PASS: TestCreateTask_TableDriven (0.00s)
=== RUN   TestGetTaskByID
--- PASS: TestGetTaskByID (0.00s)
=== RUN   TestGetTaskByID_NotFound
--- PASS: TestGetTaskByID_NotFound (0.00s)
PASS
ok      github.com/go-kruda/tutorial/04-advanced/07-testing/complete
```

Congratulations! You have written unit tests for a Kruda API using TestClient!

---

## Compare with complete/

If you get stuck, check the solution in **[complete/](./complete/)** and compare with your code:

```bash
diff starter/main.go complete/main.go
```

View the solution test file:

```bash
cat complete/handler_test.go
```

---

## Key Concepts Summary

| Concept | Description |
|---|---|
| `app.Compile()` | Compile routes before creating a test client |
| `kruda.NewTestClient(app)` | Create a test client for testing without starting a server |
| `client.Get("/path")` | Send a GET request through the test client |
| `client.Post("/path", body)` | Send a POST request with a JSON body |
| `resp.StatusCode()` | Check the HTTP status code |
| `resp.JSON(&v)` | Parse the JSON response into a struct |
| `resp.BodyString()` | Read the response body as a string |
| Table-driven tests | Standard Go pattern: define test cases as a slice of structs |
| `t.Run(name, func)` | Create a sub-test for each test case |
| `go test -v ./...` | Run all tests with verbose output |
| Fresh app per test | Create a new app in each test for isolation |
| `_test.go` convention | Go only runs tests from files ending with `_test.go` |

---

## Next Lesson

Great job! You have learned how to write unit tests for a Kruda API using TestClient. In the next lesson we will learn **Architecture** -- how to structure a project using Clean Architecture.

--> [Section 04-08 -- Architecture](../08-architecture/)
