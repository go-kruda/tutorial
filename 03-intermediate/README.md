# 🗄️ Section 03 — Intermediate: DB Integration, Docker, Config & Error Handling

⏱️ Estimated time: **45 minutes**

In this section you will learn how to connect **PostgreSQL** to Kruda, manage **configuration** via environment variables, use **Docker Compose** for development, and handle **errors** systematically with `KrudaError`.

---

## 🎯 What You Will Learn

- Connect to PostgreSQL with `database/sql` + `pgx` driver
- Read config from environment variables with fallback defaults
- Use Docker Compose to run PostgreSQL for development
- Use `c.Context()` to pass request context to database queries
- Handle errors with `kruda.NotFound()`, `kruda.InternalError()`, `kruda.BadRequest()`
- Render errors as RFC 9457 `application/problem+json` with `kruda.WithProblemJSON()`
- Enrich errors with the fluent `.WithType()`, `.WithDetail()`, `.WithInstance()`, `.With()` builders
- Use `app.MapError()` for automatic error mapping

---

## 📋 Prerequisites

| Tool | Version |
|-----------|---------|
| Go | 1.25+ |
| Docker & Docker Compose | For running PostgreSQL |
| Git | Latest |

> You should complete [Section 01 — Beginner](../01-beginner/) and [Section 02 — Auto CRUD](../02-auto-crud/) first

---

## 📁 File Structure

```
03-intermediate/
├── README.md              <-- You are here
├── docker-compose.yml     <-- PostgreSQL container
├── starter/               <-- Starter code (with TODOs to fill in)
│   ├── go.mod
│   └── main.go
└── complete/              <-- Reference implementation
    ├── go.mod
    └── main.go
```

---

## 🐘 Step 1: Start PostgreSQL with Docker Compose

```bash
cd 03-intermediate
docker compose up -d
```

Create the users table:

```bash
docker compose exec postgres psql -U postgres -d tutorial -c "
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);"
```

---

## ⚙️ Step 2: Configuration from Environment Variables

```go
type Config struct {
    Port   string
    DBHost string
    DBPort string
    DBUser string
    DBPass string
    DBName string
}

func loadConfig() Config {
    return Config{
        Port:   envOrDefault("APP_PORT", "3000"),
        DBHost: envOrDefault("DB_HOST", "localhost"),
        DBPort: envOrDefault("DB_PORT", "5432"),
        DBUser: envOrDefault("DB_USER", "postgres"),
        DBPass: envOrDefault("DB_PASS", "postgres"),
        DBName: envOrDefault("DB_NAME", "tutorial"),
    }
}

func envOrDefault(key, fallback string) string {
    v := os.Getenv(key)
    if v == "" {
        return fallback
    }
    return v
}
```

> Kruda does not enforce a config library — you can use `os.Getenv` directly or use a library like `envconfig`

---

## 🔌 Step 3: Connect to the Database

```go
import (
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"
)

func connectDB(cfg Config) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName,
    )
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, fmt.Errorf("open database: %w", err)
    }
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("ping database: %w", err)
    }
    return db, nil
}
```

> Kruda does not have a built-in database layer — it uses Go's standard `database/sql`, making it portable and giving you full control over connection pooling

---

## 📡 Step 4: Use `c.Context()` with Database Queries

Key point: always pass the request context to database queries so that queries are cancelled when the client disconnects.

```go
kruda.Get[struct{}, []UserResponse](app, "/users", func(c *kruda.C[struct{}]) (*[]UserResponse, error) {
    rows, err := db.QueryContext(c.Context(), "SELECT id, name, email FROM users ORDER BY id")
    if err != nil {
        return nil, kruda.InternalError(fmt.Sprintf("query users: %v", err))
    }
    defer rows.Close()

    var users []UserResponse
    for rows.Next() {
        var u UserResponse
        if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
            return nil, kruda.InternalError(fmt.Sprintf("scan user: %v", err))
        }
        users = append(users, u)
    }
    return &users, nil
})
```

---

## ❌ Step 5: Error Handling with KrudaError

Kruda provides built-in error helpers:

| Function | HTTP Status | Use When |
|----------|------------|---------|
| `kruda.BadRequest(msg)` | 400 | Invalid input |
| `kruda.NotFound(msg)` | 404 | Data not found |
| `kruda.InternalError(msg)` | 500 | Server error |
| `kruda.Unauthorized(msg)` | 401 | Not authenticated |
| `kruda.Forbidden(msg)` | 403 | No permission |
| `kruda.Conflict(msg)` | 409 | Duplicate data |
| `kruda.NewError(code, msg)` | custom | Custom status code |

Example usage with `sql.ErrNoRows`:

```go
kruda.Get[GetUserInput, UserResponse](app, "/users/:id", func(c *kruda.C[GetUserInput]) (*UserResponse, error) {
    var user UserResponse
    err := db.QueryRowContext(c.Context(),
        "SELECT id, name, email FROM users WHERE id = $1", c.In.ID,
    ).Scan(&user.ID, &user.Name, &user.Email)

    if err == sql.ErrNoRows {
        return nil, kruda.NotFound(fmt.Sprintf("user with id %d not found", c.In.ID))
    }
    if err != nil {
        return nil, kruda.InternalError(fmt.Sprintf("query user: %v", err))
    }
    return &user, nil
})
```

### RFC 9457 `problem+json` Errors

Since kruda **v1.4.0**, you can opt into [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457) `application/problem+json` error responses with one option:

```go
app := kruda.New(
    kruda.WithValidator(kruda.NewValidator()),
    kruda.WithProblemJSON(),
)
```

> ⚠️ Fix: `CreateUserInput` already carries `validate:"required"` / `validate:"required,email"` tags (see Step 2's types), but until now no lesson in this tutorial ever called `kruda.WithValidator(...)` -- so those tags did nothing. `kruda.WithValidator(kruda.NewValidator())` turns them on for the first time.

With `WithProblemJSON()` on, every `KrudaError` renders as a problem document instead of the plain `{code, message}` shape. Chain the fluent builders to add RFC 9457 fields:

```go
kruda.Get[GetUserInput, UserResponse](app, "/users/:id", func(c *kruda.C[GetUserInput]) (*UserResponse, error) {
    var user UserResponse
    err := db.QueryRowContext(c.Context(),
        "SELECT id, name, email FROM users WHERE id = $1", c.In.ID,
    ).Scan(&user.ID, &user.Name, &user.Email)

    if err == sql.ErrNoRows {
        return nil, kruda.NotFound(fmt.Sprintf("user with id %d not found", c.In.ID)).
            WithType("https://errors.example.com/not-found").
            With("userId", c.In.ID)
    }
    if err != nil {
        return nil, kruda.InternalError(fmt.Sprintf("query user: %v", err))
    }
    return &user, nil
})
```

`GET /users/999` now returns:

```json
{
  "type": "https://errors.example.com/not-found",
  "title": "Not Found",
  "status": 404,
  "detail": "user with id 999 not found",
  "instance": "/users/999",
  "userId": 999
}
```

| Builder | Sets |
|---|---|
| `.WithType(uri)` | RFC 9457 `type` member (defaults to `"about:blank"`) |
| `.WithDetail(text)` | `detail` member (defaults to the error's message) |
| `.WithInstance(uri)` | `instance` member (defaults to the request path) |
| `.With(key, value)` | An arbitrary extension member (reserved names `type`/`title`/`status`/`detail`/`instance`/`errors` are dropped) |

`title` always mirrors the HTTP status text and is not settable. Validation failures (from `kruda.WithValidator`) short-circuit to their own shape -- `title: "Validation failed"`, `status: 422`, and an `errors` array of per-field problems -- shown in Step 7.

> 💡 If you also call `kruda.WithErrorHandler(...)`, it takes precedence over `WithProblemJSON()` -- the two are mutually exclusive in effect, and this lesson doesn't use `WithErrorHandler`.

---

## 🔗 Step 6: Error Mapping (Optional)

`app.MapError()` lets you automatically map Go errors to HTTP responses:

```go
var ErrUserNotFound = fmt.Errorf("user not found")

app.MapError(ErrUserNotFound, 404, "user not found")
```

For type-based mapping use `kruda.MapErrorType[T]()`:

```go
kruda.MapErrorType[*ValidationError](app, 422, "validation failed")
```

---

## 🧪 Step 7: Run and Test

```bash
cd starter
go run main.go
```

Test with curl:

```bash
# Create a user
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'

# List users
curl http://localhost:3000/users

# Get a user by ID
curl http://localhost:3000/users/1

# Delete a user
curl -X DELETE http://localhost:3000/users/1

# Test 404 -- now a problem+json body with a "type" and "userId" extension
curl http://localhost:3000/users/999

# Test validation -- 422 problem+json with a field-level "errors" array
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"","email":"not-an-email"}'
```

The validation request returns:

```json
{
  "type": "about:blank",
  "title": "Validation failed",
  "status": 422,
  "detail": "validation failed: 2 errors",
  "instance": "/users",
  "errors": [
    {"field": "name", "rule": "required", "param": "", "message": "name is required", "value": ""},
    {"field": "email", "rule": "email", "param": "", "message": "email must be a valid email address", "value": "not-an-email"}
  ]
}
```

---

## 🔍 Compare with complete/

```bash
diff starter/main.go complete/main.go
```

---

## 📝 Key Concepts Summary

| Concept | Description |
|---------|-------------|
| `database/sql` + `pgx` | Connect to PostgreSQL using the standard library |
| `c.Context()` | Pass request context to DB queries (cancel on disconnect) |
| `envOrDefault()` | Read config from env vars with fallback |
| `docker compose` | Run PostgreSQL for development |
| `kruda.NotFound()` | Return 404 error |
| `kruda.InternalError()` | Return 500 error |
| `kruda.BadRequest()` | Return 400 error |
| `kruda.WithProblemJSON()` | Render errors as RFC 9457 `application/problem+json` |
| `.WithType()` / `.WithDetail()` / `.WithInstance()` / `.With()` | Fluent `*KrudaError` builders for problem+json fields and extensions |
| `kruda.WithValidator(kruda.NewValidator())` | Activate struct `validate` tags -- required to get 422s |
| `app.MapError()` | Automatically map Go error → HTTP status |
| `kruda.MapErrorType[T]()` | Map error type → HTTP status |

---

## ➡️ Next Section

Awesome! You have learned database integration, config management, and error handling. The next section moves into Advanced topics starting with **DI Container**

→ [Section 04-01 — DI Container](../04-advanced/01-di-container/)
