# 🔐 Section 04-02 — Auth Middleware: JWT Authentication with Middleware Chain

⏱️ Estimated time: **30 minutes**

Welcome to the Auth Middleware lesson! 🎉 In this lesson you'll learn how to build a **JWT Authentication Middleware** using Kruda's **Middleware Chain** — protect your API endpoints with automatic token verification without having to write auth logic repeatedly in every handler.

---

## 🎯 Lesson Objectives

- Understand the concept of **Middleware Chain** and why it matters
- Learn how to create middleware with `kruda.HandlerFunc`
- Build a JWT authentication middleware that validates Bearer tokens
- Use `app.Group()` and `api.Guard()` to separate public/protected routes
- Store user data in context with `c.Set()` / `c.Get()`

---

## 📚 What You'll Learn

By the end of this lesson you'll be able to:

- ✅ Create a middleware function with `kruda.HandlerFunc`
- ✅ Validate JWT tokens from the `Authorization` header
- ✅ Use `c.Next()` to pass the request to the next handler
- ✅ Use `app.Group()` to group routes
- ✅ Use `api.Guard()` to apply middleware to a group of routes
- ✅ Store and retrieve data from context with `c.Set()` / `c.Get()`
- ✅ Separate public routes (e.g. `/login`) from protected routes (e.g. `/api/profile`)

---

## 📋 Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25 or higher |
| Git | Latest version |
| Text Editor / IDE | VS Code, GoLand, or your preferred editor |
| Terminal | For running `go run` |

> 💡 If you haven't completed Section 04-01 (DI Container) yet, it's recommended to go back and do it first — 👉 [Section 04-01 — DI Container](../01-di-container/)

---

## 📂 File Structure

```
04-advanced/02-auth-middleware/
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

## 💡 Why Use Middleware Chain?

Imagine you have 20 endpoints that need JWT token verification — if you write auth logic in every handler, you'll run into problems:

| Problem | How Middleware Chain Helps |
|---|---|
| Auth code duplicated in every handler | Write once, apply to all routes |
| Forgetting to add auth check in a new handler | Middleware runs automatically for all routes in the group |
| Difficult to change auth logic | Fix it in one place, all routes get the update immediately |
| Can't test auth separately from handlers | Middleware is a separate function that can be tested independently |

> 🛡️ Kruda's Middleware Chain works as a **chain of responsibility** — the request passes through each middleware one by one before reaching the handler

```
Request → Middleware 1 → Middleware 2 → Handler → Response
```

---

## 🛠️ Step-by-Step Guide

### Step 1: Open the starter project

```bash
cd 04-advanced/02-auth-middleware/starter
```

Open the `main.go` file — you'll see a structure with `// TODO:` comments indicating where you need to add code.

### Step 2: Understand the structure

In this lesson we'll build an API with 2 parts:

```
Public Routes:     POST /login, GET /health
Protected Routes:  GET /api/profile (requires JWT token)
```

- **JWT Helpers** — Generate and validate tokens (provided for you)
- **Auth Middleware** — Validate the Authorization header
- **Typed Route Functions** — `kruda.Post[In, Out]()` / `kruda.Get[In, Out]()` / `kruda.GroupGet[In, Out]()` for each endpoint

### Step 3: Write the Auth Middleware

Replace the `// TODO:` in `authMiddleware()`:

```go
func authMiddleware() kruda.HandlerFunc {
    return func(c *kruda.Ctx) error {
        // 1. Read the Authorization header
        authHeader := c.Header("Authorization")
        if authHeader == "" {
            return kruda.Unauthorized("missing Authorization header")
        }

        // 2. Validate the "Bearer <token>" format
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            return kruda.Unauthorized("invalid Authorization format (expected: Bearer <token>)")
        }

        // 3. Validate token
        claims, err := validateToken(parts[1])
        if err != nil {
            return kruda.Unauthorized(fmt.Sprintf("authentication failed: %v", err))
        }

        // 4. Store user data in context
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)

        // 5. Pass to the next handler
        return c.Next()
    }
}
```

> 🔑 `c.Next()` is the heart of the Middleware Chain — if you don't call `c.Next()`, the request will stop at this middleware and never reach the handler

### Step 4: Register Routes in `main()`

Replace the `// TODO:` in `main()` — register public routes, create a protected group, and write inline handlers with `kruda.Post` / `kruda.Get` / `kruda.GroupGet`:

```go
func main() {
    app := kruda.New()

    // ── Public routes — no auth required ──────────────────────
    kruda.Get[struct{}, MessageResponse](app, "/health", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
        return &MessageResponse{Message: "OK"}, nil
    })

    kruda.Post[LoginInput, TokenResponse](app, "/login", func(c *kruda.C[LoginInput]) (*TokenResponse, error) {
        if c.In.Username == "admin" && c.In.Password == "secret" {
            token, err := generateToken(c.In.Username, "admin")
            if err != nil {
                return nil, kruda.InternalError("generate token failed")
            }
            return &TokenResponse{Token: token, ExpiresIn: 3600}, nil
        }
        return nil, kruda.Unauthorized("invalid username or password")
    })

    // ── Protected routes — JWT token required ───────────────
    api := app.Group("/api")
    api.Guard(authMiddleware())

    kruda.GroupGet[struct{}, ProfileResponse](api, "/profile", func(c *kruda.C[struct{}]) (*ProfileResponse, error) {
        username := c.Get("username").(string)
        role := c.Get("role").(string)
        return &ProfileResponse{
            Username: username,
            Role:     role,
            Message:  fmt.Sprintf("Welcome back, %s! You have %s access.", username, role),
        }, nil
    })

    log.Println("Server starting on :3000 ...")
    log.Fatal(app.Listen(":3000"))
}
```

> ✨ Kruda automatically deserialises the JSON body into `LoginInput` via `c.In` — no need to call `json.NewDecoder` yourself!

> 🧩 Notice that the profile handler doesn't need to verify the token itself — the middleware already did that! The handler simply retrieves data from the context that the middleware prepared with `c.Get("username").(string)`

> 🛡️ `app.Group("/api")` creates a route group with the `/api` prefix, and `api.Guard(authMiddleware())` applies the auth middleware to all routes in this group — `Guard()` is a semantic alias for `Use()` that reads more clearly for auth/permission middleware

### Step 5: Run and test

```bash
go run main.go
```

Open another terminal window and test with `curl`:

```bash
# 1. Test health check (public)
curl http://localhost:3000/health

# 2. Try accessing a protected route without a token (will get 401)
curl http://localhost:3000/api/profile

# 3. Login to get a token
curl -X POST http://localhost:3000/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"secret"}'

# 4. Use the token to access the protected route
curl http://localhost:3000/api/profile \
  -H "Authorization: Bearer <token-from-step-3>"
```

If everything works correctly:
- Step 1: Returns `{"message":"OK"}`
- Step 2: Returns `{"message":"missing Authorization header"}` (401)
- Step 3: Returns `{"token":"eyJ...","expires_in":3600}`
- Step 4: Returns `{"username":"admin","role":"admin","message":"Welcome back, admin! You have admin access."}` 🎉

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
| `kruda.HandlerFunc` | A middleware function that takes `*kruda.Ctx` and returns `error` |
| `c.Next()` | Passes the request to the next middleware/handler in the chain |
| `c.Header("Authorization")` | Reads an HTTP header value from the request |
| `kruda.Unauthorized()` | Creates an error response with HTTP 401 |
| `app.Group("/api")` | Creates a route group with a shared prefix |
| `api.Guard()` | Applies auth middleware to all routes in the group (semantic alias for `Use()`) |
| `c.Set()` / `c.Get()` | Store/retrieve data in the request context |
| `kruda.Post[In, Out]()` | Register a typed POST handler on the app |
| `kruda.GroupGet[In, Out]()` | Register a typed GET handler on a group |
| Bearer Token | The `Authorization: Bearer <token>` format for JWT |
| Public vs Protected | Separating routes that don't require auth from routes that do |

---

## ➡️ Next Lesson

Awesome! 🎊 You've learned how to build JWT Authentication Middleware using Kruda's Middleware Chain. In the next lesson we'll learn about the **OpenAPI Generator** — how to automatically generate API documentation from Typed Handlers 📄

👉 [Section 04-03 — OpenAPI](../03-openapi/)
