package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-kruda/kruda"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// ============================================================
// Application Configuration
// ============================================================

// Config holds all application settings loaded from the
// environment.
type Config struct {
	Port   string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

// loadConfig reads configuration from environment variables
// and returns a Config struct. Missing variables fall back to
// sensible defaults for local development.
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

// ============================================================
// Database Connection
// ============================================================
//
// Kruda does not include a built-in database layer. Use Go's
// standard database/sql package with a driver like pgx. This
// keeps your application portable and gives you full control
// over connection pooling, transactions, and query patterns.

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

// ============================================================
// Request / Response Types
// ============================================================

// CreateUserInput represents the JSON body for creating a user.
type CreateUserInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// GetUserInput captures the :id path parameter.
type GetUserInput struct {
	ID int `param:"id"`
}

// UserResponse represents the JSON response returned for a user.
type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// MessageResponse is a generic message envelope.
type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	// ── 1. Load Configuration ─────────────────────────────────
	cfg := loadConfig()

	// ── 2. Connect to the Database ────────────────────────────
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// ── 3. Create the Kruda Application ───────────────────────
	app := kruda.New()

	// ── 4. Register Routes ────────────────────────────────────
	//
	// Each handler uses the db handle via closure -- this is
	// idiomatic Go dependency injection without a framework.

	// GET /users -- list all users
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

	// POST /users -- create a new user
	kruda.Post[CreateUserInput, UserResponse](app, "/users", func(c *kruda.C[CreateUserInput]) (*UserResponse, error) {
		var user UserResponse
		err := db.QueryRowContext(c.Context(),
			"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email",
			c.In.Name, c.In.Email,
		).Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, kruda.InternalError(fmt.Sprintf("create user: %v", err))
		}
		return &user, nil
	})

	// GET /users/:id -- get a single user
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

	// DELETE /users/:id -- delete a user
	kruda.Delete[GetUserInput, MessageResponse](app, "/users/:id", func(c *kruda.C[GetUserInput]) (*MessageResponse, error) {
		result, err := db.ExecContext(c.Context(),
			"DELETE FROM users WHERE id = $1", c.In.ID,
		)
		if err != nil {
			return nil, kruda.InternalError(fmt.Sprintf("delete user: %v", err))
		}

		n, _ := result.RowsAffected()
		if n == 0 {
			return nil, kruda.NotFound(fmt.Sprintf("user with id %d not found", c.In.ID))
		}

		return &MessageResponse{
			Message: fmt.Sprintf("user %d deleted successfully", c.In.ID),
		}, nil
	})

	// ── 5. Start the Server ───────────────────────────────────
	log.Printf("Server starting on :%s ...\n", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
