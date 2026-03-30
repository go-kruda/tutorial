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

# Test 404
curl http://localhost:3000/users/999
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
| `app.MapError()` | Automatically map Go error → HTTP status |
| `kruda.MapErrorType[T]()` | Map error type → HTTP status |

---

## ➡️ Next Section

Awesome! You have learned database integration, config management, and error handling. The next section moves into Advanced topics starting with **DI Container**

→ [Section 04-01 — DI Container](../04-advanced/01-di-container/)
