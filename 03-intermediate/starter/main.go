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

type Config struct {
	Port   string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

// TODO: implement loadConfig -- read config from environment variables
//
// Hint: Use os.Getenv() with fallback defaults
//   Port:   envOrDefault("APP_PORT", "3000")
//   DBHost: envOrDefault("DB_HOST", "localhost")
//   etc.
func loadConfig() Config {
	return Config{}
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

// TODO: implement connectDB -- connect to PostgreSQL using database/sql
//
// Hint: Use sql.Open("pgx", dsn) then db.Ping() to verify the connection
//   dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", ...)
func connectDB(cfg Config) (*sql.DB, error) {
	return nil, fmt.Errorf("not implemented")
}

// ============================================================
// Request / Response Types
// ============================================================

type CreateUserInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type GetUserInput struct {
	ID int `param:"id"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	cfg := loadConfig()

	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	app := kruda.New()

	// TODO: Register route GET /users -- fetch all users
	//
	// Hint: Use db.QueryContext(c.Context(), "SELECT ...") then scan the results
	//
	//   kruda.Get[struct{}, []UserResponse](app, "/users", func(c *kruda.C[struct{}]) (*[]UserResponse, error) {
	//       rows, err := db.QueryContext(c.Context(), "SELECT id, name, email FROM users ORDER BY id")
	//       // scan rows into []UserResponse
	//       return &users, nil
	//   })

	// TODO: Register route POST /users -- create a new user
	//
	// Hint: Use db.QueryRowContext + RETURNING to get the ID back
	//
	//   kruda.Post[CreateUserInput, UserResponse](app, "/users", func(c *kruda.C[CreateUserInput]) (*UserResponse, error) {
	//       var user UserResponse
	//       err := db.QueryRowContext(c.Context(),
	//           "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email",
	//           c.In.Name, c.In.Email,
	//       ).Scan(&user.ID, &user.Name, &user.Email)
	//       return &user, nil
	//   })

	// TODO: Register route GET /users/:id -- fetch a user by ID
	//
	// Hint: c.In.ID will be automatically parsed from :id (param:"id" tag)
	//   Use sql.ErrNoRows to check if no data was found
	//   return nil, kruda.NotFound("user not found")

	// TODO: Register route DELETE /users/:id -- delete a user by ID
	//
	// Hint: Use db.ExecContext + RowsAffected() to verify the deletion succeeded

	_ = db
	log.Printf("Server starting on :%s ...\n", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
