# 🏗️ Section 04-08 — Architecture: Clean Architecture with DI Container

⏱️ Estimated time: **45 minutes**

Welcome to the final lesson of the Advanced Section! In this lesson you will learn how to structure a project using **Clean Architecture** by splitting code into 3 layers -- handler, service, repository -- each layer in its own package, and using Kruda's **DI Container** to wire all dependencies together.

---

## Learning Objectives

- Understand the principles of **Clean Architecture** and why it matters for large projects
- Split code into 3 layers: **Handler**, **Service**, **Repository**
- Use the **DI Container** as a composition root to wire dependencies
- See how the dependency direction flows from outside in (handler -> service -> repository)
- Understand that each layer only knows about the layer directly below it

---

## What You Will Learn

By the end of this lesson you will be able to:

- Structure a project with multiple packages (handler/, service/, repository/)
- Create a **Repository layer** for data access
- Create a **Service layer** for business logic and validation
- Create a **Handler layer** that registers routes with `kruda.Get`, `kruda.Post`, `kruda.Delete`
- Use `kruda.NewContainer()`, `container.Give()`, `kruda.MustUse[T]()` to wire all layers together
- Understand the composition root pattern in `main()`

---

## Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25+ |
| Git | Latest |
| Text Editor / IDE | VS Code, GoLand, or your preferred editor |
| Terminal | For running `go run` |

> This lesson builds on Section 04-01 (DI Container) -- if you haven't completed it, it's recommended to go back and do it first -- [Section 04-01 -- DI Container](../01-di-container/)

---

## File Structure

```
04-advanced/08-architecture/
├── README.md              <- You are here
├── starter/               <- Starter code (single-file, with TODOs to fill in)
│   ├── go.mod
│   └── main.go
└── complete/              <- Complete solution with multi-package layout
    ├── go.mod
    ├── main.go            <- Composition root (DI wiring)
    ├── handler/
    │   └── user.go        <- Route registration (RegisterRoutes)
    ├── service/
    │   └── user.go        <- Business logic & validation
    └── repository/
        └── user.go        <- Data access (in-memory store)
```

- **[starter/](./starter/)** -- All code in a single file (`main.go`) with `// TODO:` markers for you to complete -- compiles immediately
- **[complete/](./complete/)** -- Full working solution split into packages following clean architecture

---

## Why Clean Architecture?

As a project grows, writing everything in a single file leads to:

| Problem | How Clean Architecture Helps |
|---|---|
| Long code that's hard to navigate | Split into packages by responsibility |
| Business logic coupled to HTTP | Service layer separated from handler |
| Hard to change databases | Repository layer abstracts data access |
| Hard to test | Each layer can be tested independently |
| Tangled dependencies | DI Container manages the dependency graph |

### Dependency Direction

```
+---------------------------------------------+
|                  main.go                     |
|            (Composition Root)                |
|         DI Container wires all layers        |
+---------------------+------------------------+
                      |
         +------------+------------+
         v            v            v
   +----------+  +--------+  +--------------+
   | handler/ |->|service/|->| repository/  |
   | (HTTP)   |  |(Logic) |  | (Data)       |
   +----------+  +--------+  +--------------+
```

> Arrows point inward -- handler knows about service but service doesn't know about handler; service knows about repository but repository doesn't know about service.

---

## Step-by-Step Guide

### Step 1: Open the starter project

```bash
cd 04-advanced/08-architecture/starter
```

Open the `main.go` file -- you will see all 3 layers in a single file with `// TODO:` comments.

### Step 2: Understand the 3 Layers

The starter has 3 main parts:

1. **Repository** (`UserRepository`) -- Manages user data in memory
2. **Service** (`UserService`) -- Business logic such as validation
3. **Route Registration** -- The `registerRoutes()` function that registers routes on the app

### Step 3: Implement Service Validation

Replace the `// TODO:` in `CreateUser()` of UserService:

```go
func (s *UserService) CreateUser(name, email string) (User, error) {
    if name == "" {
        return User{}, fmt.Errorf("name is required")
    }
    if email == "" {
        return User{}, fmt.Errorf("email is required")
    }
    return s.repo.Create(name, email), nil
}
```

> Business rules live in the service layer -- the handler doesn't need to know about validation.

### Step 4: Implement Route Handlers

Replace the `// TODO:` in `registerRoutes()` to register all 4 routes. For example, GET /users:

```go
kruda.Get[struct{}, []User](app, "/users", func(c *kruda.C[struct{}]) (*[]User, error) {
    users := svc.ListUsers()
    return &users, nil
})
```

For routes that need a path parameter, use an input struct with the `param` tag:

```go
kruda.Get[GetUserInput, User](app, "/users/:id", func(c *kruda.C[GetUserInput]) (*User, error) {
    user, err := svc.GetUser(c.In.ID)
    if err != nil {
        return nil, kruda.NotFound(err.Error())
    }
    return &user, nil
})
```

> The handler is just a thin adapter -- receive the request, call the service, return the response.

### Step 5: Wire Dependencies with DI Container

Replace the `// TODO:` in `main()`:

```go
container := kruda.NewContainer()

repo := NewUserRepository()
container.Give(repo)

svc := NewUserService(repo)
container.Give(svc)

userService := kruda.MustUse[*UserService](container)

app := kruda.New()
registerRoutes(app, userService)

log.Fatal(app.Listen(":3000"))
```

> `main()` serves as the **composition root** -- the only place that knows about every layer and wires them together.

### Step 6: (Challenge) Split into Packages

Once the starter is working, challenge yourself by splitting the code into packages:

```
starter/
├── main.go            <- composition root
├── handler/
│   └── user.go        <- Move registerRoutes + input types here
├── service/
│   └── user.go        <- Move UserService here
└── repository/
    └── user.go        <- Move UserRepository here
```

> Check the solution in **[complete/](./complete/)** to compare what each package should contain.

### Step 7: Run and Test

```bash
go run main.go
```

Test with `curl`:

```bash
# Create a new user
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"somchai","email":"somchai@example.com"}'

# List all users
curl http://localhost:3000/users

# Get a user by ID
curl http://localhost:3000/users/1

# Delete a user
curl -X DELETE http://localhost:3000/users/1
```

---

## Compare with complete/

If you get stuck, check the solution in **[complete/](./complete/)** -- notice how the code is split into packages:

```bash
# View the structure
ls -la complete/
ls -la complete/handler/
ls -la complete/service/
ls -la complete/repository/
```

---

## Key Concepts Summary

| Concept | Description |
|---|---|
| Clean Architecture | Split code into layers by responsibility |
| Repository layer | Data access -- only knows how to store/retrieve data |
| Service layer | Business logic -- validation, orchestration |
| Handler layer | Route registration -- registers routes with `kruda.Get`, `kruda.Post`, etc. |
| `RegisterRoutes(app, svc)` | Function in the handler package that registers all routes |
| `kruda.C[T]` | Context with input (JSON body, path params) automatically bound |
| `kruda.BadRequest()` / `kruda.NotFound()` | Error helpers for HTTP error responses |
| Composition root | `main()` that wires all layers with DI Container |
| Dependency direction | handler -> service -> repository (outside in) |
| `kruda.NewContainer()` | Create a DI Container |
| `container.Give()` | Register an instance (singleton) into the container |
| `kruda.MustUse[T]()` | Retrieve a service from the container |

---

## Next Lesson

Excellent work! You have completed the entire Advanced Section! You now have the full knowledge to build a well-structured, testable Kruda application ready for production. In the next section we will move into **Section Production** -- learning about monitoring, deployment, and benchmarking.

--> [Section 05 -- Production](../../05-production/)
