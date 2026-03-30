package main

import (
	"log"

	"github.com/go-kruda/kruda"
	"github.com/go-kruda/kruda/middleware"
)

// ============================================================
// Request / Response Types
// ============================================================

// CreateBookInput represents the JSON body for creating a book.
type CreateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

// GetBookByIDInput captures the :id path parameter.
// tag `param:"id"` tells Kruda to parse :id from URL automatically.
type GetBookByIDInput struct {
	ID int `param:"id"`
}

// BookResponse represents the JSON response returned for a book.
type BookResponse struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// MessageResponse is a generic message envelope.
type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	app := kruda.New()

	// TODO: Add Recovery + Logger middleware
	//
	// Hint: app.Use(middleware.Recovery(), middleware.Logger())
	app.Use(middleware.Recovery(), middleware.Logger())

	// TODO: Register route GET /books -- fetch all books
	//
	// Hint: Use kruda.Get[In, Out] where In is struct{} (no input body)
	//
	//   kruda.Get[struct{}, []BookResponse](app, "/books", func(c *kruda.C[struct{}]) (*[]BookResponse, error) {
	//       result := make([]BookResponse, len(books))
	//       copy(result, books)
	//       return &result, nil
	//   })

	// TODO: Register route POST /books -- create a new book
	//
	// Hint: Use kruda.Post[CreateBookInput, BookResponse]
	//   Kruda will automatically parse the JSON body into CreateBookInput
	//   Access data via c.In.Title, c.In.Author
	//
	//   kruda.Post[CreateBookInput, BookResponse](app, "/books", func(c *kruda.C[CreateBookInput]) (*BookResponse, error) {
	//       book := BookResponse{Title: c.In.Title, Author: c.In.Author}
	//       return &book, nil
	//   })

	// TODO: Register route GET /books/:id -- fetch a book by ID
	//
	// Hint: Use GetBookByIDInput which has the `param:"id"` tag
	//   Kruda will automatically parse :id from the URL as int via c.In.ID
	//   Use kruda.NotFound("msg") when data is not found
	//
	//   kruda.Get[GetBookByIDInput, BookResponse](app, "/books/:id", func(c *kruda.C[GetBookByIDInput]) (*BookResponse, error) {
	//       // Use c.In.ID to find the book
	//       return nil, kruda.NotFound("not found")
	//   })

	// TODO: Register route DELETE /books/:id -- delete a book by ID
	//
	// Hint: Use kruda.Delete[GetBookByIDInput, MessageResponse]
	//
	//   kruda.Delete[GetBookByIDInput, MessageResponse](app, "/books/:id", func(c *kruda.C[GetBookByIDInput]) (*MessageResponse, error) {
	//       return nil, kruda.NotFound("not found")
	//   })

	// TODO: Start the server on port 3000
	//
	// Hint: log.Fatal(app.Listen(":3000"))
	_ = app
	log.Println("Server starting on :3000 ...")
}
