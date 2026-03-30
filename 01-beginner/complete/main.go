package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-kruda/kruda"
	"github.com/go-kruda/kruda/middleware"
)

// ============================================================
// Request / Response Types
// ============================================================
//
// These structs define the shape of JSON payloads flowing in and
// out of our API. Kruda's Typed Handler system uses Go generics
// to bind these types at compile time -- the framework handles
// JSON marshalling/unmarshalling automatically, so you never
// call json.NewEncoder or json.NewDecoder yourself.

// CreateBookInput represents the JSON body for creating a book.
type CreateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

// GetBookByIDInput captures the :id path parameter.
// The `param:"id"` tag tells Kruda to auto-parse the :id segment
// from the URL and convert it to an int -- no manual strconv needed.
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
// In-Memory Book Store
// ============================================================
//
// For this beginner tutorial we keep things simple with an
// in-memory slice. A mutex protects concurrent access so the
// server is safe under parallel requests.

var (
	books  []BookResponse
	mu     sync.Mutex
	nextID = 1
)

// ============================================================
// Application Entry Point
// ============================================================
//
// Kruda's Typed Handlers leverage Go generics to give you:
//
//  1. Compile-time type safety -- if your handler returns the
//     wrong type, the compiler catches it immediately.
//  2. Automatic JSON deserialisation -- the framework reads the
//     request body and populates your input struct for you.
//  3. Automatic JSON serialisation -- just return your response
//     value and Kruda writes the JSON with correct headers.
//  4. Automatic path parameter parsing -- use `param:"id"` tags
//     and Kruda extracts and converts URL segments for you.
//
// Handler registration pattern:
//
//   kruda.Get[InputType, OutputType](app, "/path", handlerFunc)
//
// The handler function signature is:
//
//   func(c *kruda.C[InputType]) (*OutputType, error)
//
// Access the parsed input via c.In -- for path params, query
// params, or JSON body, all bound from struct tags.

func main() {
	// Create a new Kruda application instance.
	app := kruda.New()

	// Add essential middleware: Recovery catches panics and
	// returns 500 instead of crashing; Logger logs every
	// request with method, path, status, and latency.
	app.Use(middleware.Recovery(), middleware.Logger())

	// ── GET /books -- list all books ─────────────────────────
	//
	// The input type is struct{} because GET /books has no
	// request body or path parameters. The output type is
	// []BookResponse -- Kruda serialises the slice to a JSON
	// array automatically.
	kruda.Get[struct{}, []BookResponse](app, "/books", func(c *kruda.C[struct{}]) (*[]BookResponse, error) {
		mu.Lock()
		defer mu.Unlock()

		// Return a copy so callers cannot mutate our store.
		result := make([]BookResponse, len(books))
		copy(result, books)
		return &result, nil
	})

	// ── POST /books -- create a new book ─────────────────────
	//
	// Kruda automatically deserialises the incoming JSON into
	// CreateBookInput -- no need to call json.NewDecoder.
	// Access the parsed fields via c.In.Title and c.In.Author.
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

		// Kruda serialises this BookResponse to JSON and sets
		// Content-Type: application/json for us.
		return &book, nil
	})

	// ── GET /books/:id -- get a single book ──────────────────
	//
	// GetBookByIDInput has `ID int` with `param:"id"` tag.
	// Kruda parses the :id URL segment and converts it to an
	// int automatically -- no strconv.Atoi needed.
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

	// ── DELETE /books/:id -- delete a book ───────────────────
	//
	// Same param:"id" pattern for path parameter extraction.
	// Returns a MessageResponse confirming the deletion.
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

	// Start the server on port 3000.
	log.Println("Server starting on :3000 ...")
	log.Fatal(app.Listen(":3000"))
}
