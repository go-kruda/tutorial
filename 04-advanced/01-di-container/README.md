# 🏗️ Section 04-01 — DI Container: Managing Dependencies Systematically

⏱️ Estimated time: **30 minutes**

Welcome to the first lesson of the Advanced Section! 🎉 In this lesson you will learn how to use the **DI Container** (Dependency Injection Container) that comes built into Kruda — no need to install any external libraries. It helps you manage dependencies between services systematically, reduces hard-coding, and makes your code easier to test.

---

## 🎯 Lesson Objectives

- Understand the concept of **Dependency Injection** and why it matters
- Learn how to create a DI Container with `kruda.NewContainer()`
- Register singleton instances with `container.Give()`
- Retrieve services with `kruda.MustUse[T]()`
- See how the container manages the dependency graph automatically

---

## 📚 What You Will Learn

By the end of this lesson you will be able to:

- ✅ Create a DI Container with `kruda.NewContainer()`
- ✅ Register a singleton instance with `container.Give()`
- ✅ Resolve a service from the container with `kruda.MustUse[T]()`
- ✅ Separate **Repository** (data access) from **Service** (business logic)
- ✅ Use the DI Container together with `kruda.Get[]` / `kruda.Post[]` / `kruda.Delete[]` to build a well-structured REST API
- ✅ Understand how the container manages the order of dependency creation automatically

---

## 📋 Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25 or higher |
| Git | Latest version |
| Text Editor / IDE | VS Code, GoLand, or your preferred editor |
| Terminal | For running `go run` |

> 💡 If you haven't completed Section 03 yet, we recommend going back to do it first since this section builds on the concepts of Typed Handlers and closure-based DI — 👉 [Section 03 — Intermediate](../../03-intermediate/)

---

## 📂 File Structure

```
04-advanced/01-di-container/
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

## 💡 Why Use a DI Container?

In Section 03 we used **closure-based dependency injection** — passing dependencies into handlers via function parameters. This works well for small apps, but as the app grows:

| Problem | How DI Container Helps |
|---|---|
| Dependencies nested multiple layers deep | Container manages the creation order automatically |
| `main()` becomes very long from wiring dependencies | Register once, resolve anywhere |
| Hard to swap implementations | Change the instance given to the container in one place, and every dependent service gets the new one |
| Hard to test | Swap implementations in the container for testing easily |

> 🧩 Kruda has a built-in DI Container — no need to import external libraries like `wire` or `dig`

---

## 🛠️ Step-by-Step

### Step 1: Open the starter project

```bash
cd 04-advanced/01-di-container/starter
```

Open the `main.go` file — you will see a structure with `// TODO:` comments indicating where you need to add code.

### Step 2: Understand the Architecture

In this lesson we will build a User API with 3 layers:

```
Handler (HTTP) → Service (Business Logic) → Repository (Data Access)
```

- **UserRepository** — manages user data (in-memory store)
- **UserService** — business logic that depends on UserRepository
- **Typed Handlers** — HTTP layer that depends on UserService

The DI Container will be the glue that connects these layers together.

### Step 3: Create the DI Container

Replace the first `// TODO:` in `main()`:

```go
container := kruda.NewContainer()
```

> 🏗️ `kruda.NewContainer()` creates an empty container — ready to accept service registrations

### Step 4: Register UserRepository

```go
repo := NewUserRepository()
container.Give(repo)
```

> 📦 `container.Give()` takes an already-created instance (singleton) and registers it in the container by type — when someone calls `MustUse`, they will get the same instance back

### Step 5: Register UserService

```go
svc := NewUserService(repo)
container.Give(svc)
```

> 🔗 Create UserService by passing the repository as a dependency, then register it in the container — every service that needs `*UserService` will get the same instance

### Step 6: Resolve Service and Create the App

**Pattern A — Resolve at startup (simple):**

```go
userService := kruda.MustUse[*UserService](container)
app := kruda.New()
```

**Pattern B — Attach container then resolve per-request (recommended):**

```go
app := kruda.New(kruda.WithContainer(container))
```

> 🎯 Pattern B is idiomatic Kruda DI — use `kruda.WithContainer(container)` when creating the app, then resolve in handlers with `kruda.MustResolve[T](c.Ctx)`. This keeps handlers decoupled from closures in main() and makes testing easier

### Step 7: Register Routes with Typed Handlers

Use `kruda.MustResolve[T](c.Ctx)` to resolve the service from the request context:

```go
kruda.Get[struct{}, []UserResponse](app, "/users", func(c *kruda.C[struct{}]) (*[]UserResponse, error) {
    svc := kruda.MustResolve[*UserService](c.Ctx)
    users := svc.ListUsers()
    return &users, nil
})

kruda.Post[CreateUserInput, UserResponse](app, "/users", func(c *kruda.C[CreateUserInput]) (*UserResponse, error) {
    svc := kruda.MustResolve[*UserService](c.Ctx)
    user, err := svc.CreateUser(c.In.Name, c.In.Email)
    if err != nil {
        return nil, kruda.BadRequest(err.Error())
    }
    return &user, nil
})

kruda.Get[GetUserInput, UserResponse](app, "/users/:id", func(c *kruda.C[GetUserInput]) (*UserResponse, error) {
    user, err := userService.GetUser(c.In.ID)
    if err != nil {
        return nil, kruda.NotFound(err.Error())
    }
    return &user, nil
})

kruda.Delete[GetUserInput, MessageResponse](app, "/users/:id", func(c *kruda.C[GetUserInput]) (*MessageResponse, error) {
    if err := userService.DeleteUser(c.In.ID); err != nil {
        return nil, kruda.NotFound(err.Error())
    }
    return &MessageResponse{
        Message: fmt.Sprintf("user %d deleted", c.In.ID),
    }, nil
})
```

> ✨ `kruda.Get[In, Out]` / `kruda.Post[In, Out]` take type parameters for input and output — the handler function receives `*kruda.C[In]` which has a field `c.In` that is automatically decoded from the request. Use `kruda.BadRequest()` / `kruda.NotFound()` for error responses

### Step 8: Run and Test

```bash
go run main.go
```

Open another terminal window and test with `curl`:

```bash
# Create a new user
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Somchai","email":"somchai@example.com"}'

# List all users
curl http://localhost:3000/users

# Get a user by ID
curl http://localhost:3000/users/1

# Delete a user
curl -X DELETE http://localhost:3000/users/1
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
| `kruda.NewContainer()` | Creates an empty DI Container ready to accept service registrations |
| `container.Give()` | Registers a singleton instance in the container by type |
| `kruda.MustUse[T]()` | Retrieves a service from the container at startup (Pattern A) |
| `kruda.WithContainer(c)` | Attaches the container to the App (Pattern B — recommended) |
| `kruda.MustResolve[T](ctx)` | Retrieves a service from the request context (Pattern B — recommended) |
| `kruda.Get[In, Out]()` | Registers a type-safe GET route with auto-decoded input |
| `kruda.Post[In, Out]()` | Registers a type-safe POST route with auto-decoded input |
| `kruda.Delete[In, Out]()` | Registers a type-safe DELETE route with auto-decoded input |
| `*kruda.C[T]` | Typed context with field `c.In` for decoded request input |
| Repository pattern | Separates data access from business logic |
| Service pattern | Centralises business rules in one place, decoupled from the HTTP layer |

---

## ➡️ Next Section

Awesome! 🎊 You have learned how to use Kruda's DI Container to manage dependencies systematically. In the next lesson we will learn about **Auth Middleware** — how to build JWT authentication middleware with Kruda's Middleware Chain 🔐

👉 [Section 04-02 — Auth Middleware](../02-auth-middleware/)
