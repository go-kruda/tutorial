package main

import (
	"log"
	"sync"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Request / Response Types
// ============================================================
//
// These structs define the shape of JSON payloads flowing in and
// out of our API. Kruda uses Go generics to bind these types at
// compile time -- the framework handles JSON marshalling and
// unmarshalling automatically.

// CreateTaskInput represents the JSON body for creating a task.
type CreateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// GetTaskInput captures the :id path parameter.
type GetTaskInput struct {
	ID int `param:"id"`
}

// TaskResponse represents the JSON response returned for a task.
type TaskResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// MessageResponse is a generic message envelope.
type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// In-Memory Task Store
// ============================================================
//
// TaskStore is a simple in-memory store that holds tasks. It is
// designed to be easily testable -- you can create a fresh store
// in each test without any external dependencies.

// TaskStore holds tasks in memory with thread-safe access.
type TaskStore struct {
	mu     sync.Mutex
	tasks  []TaskResponse
	nextID int
}

// NewTaskStore creates an empty TaskStore.
func NewTaskStore() *TaskStore {
	return &TaskStore{nextID: 1}
}

// All returns a copy of all tasks.
func (s *TaskStore) All() []TaskResponse {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]TaskResponse, len(s.tasks))
	copy(result, s.tasks)
	return result
}

// Create adds a new task and returns it.
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

// FindByID returns a task by ID, or false if not found.
func (s *TaskStore) FindByID(id int) (TaskResponse, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, t := range s.tasks {
		if t.ID == id {
			return t, true
		}
	}
	return TaskResponse{}, false
}

// ToggleDone flips the Done status of a task by ID.
// Returns the updated task and true, or zero-value and false
// if the task was not found.
func (s *TaskStore) ToggleDone(id int) (TaskResponse, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks[i].Done = !t.Done
			return s.tasks[i], true
		}
	}
	return TaskResponse{}, false
}

// ============================================================
// Route Registration
// ============================================================
//
// Why test via TestClient?
// ------------------------
// Kruda provides kruda.NewTestClient(app) so you can test your
// entire HTTP pipeline -- routing, input binding, handler logic,
// and JSON serialisation -- without starting a real HTTP server.
// You build the app, compile it, create a test client, and make
// HTTP-like requests that return a response you can assert on.

// setupApp creates and configures a Kruda app with all routes.
// This function is shared between main() and the test file so
// the test exercises the exact same routing and middleware.
func setupApp() *kruda.App {
	store := NewTaskStore()
	return setupAppWithStore(store)
}

// setupAppWithStore creates a Kruda app using the provided store.
// Tests can pass in a pre-populated store when needed.
func setupAppWithStore(store *TaskStore) *kruda.App {
	app := kruda.New()

	// GET /tasks -- list all tasks
	kruda.Get[struct{}, []TaskResponse](app, "/tasks", func(c *kruda.C[struct{}]) (*[]TaskResponse, error) {
		tasks := store.All()
		return &tasks, nil
	})

	// POST /tasks -- create a new task
	kruda.Post[CreateTaskInput, TaskResponse](app, "/tasks", func(c *kruda.C[CreateTaskInput]) (*TaskResponse, error) {
		if c.In.Title == "" {
			return nil, kruda.BadRequest("title is required")
		}
		task := store.Create(c.In.Title, c.In.Description)
		return &task, nil
	})

	// GET /tasks/:id -- get a single task
	kruda.Get[GetTaskInput, TaskResponse](app, "/tasks/:id", func(c *kruda.C[GetTaskInput]) (*TaskResponse, error) {
		task, ok := store.FindByID(c.In.ID)
		if !ok {
			return nil, kruda.NotFound("task not found")
		}
		return &task, nil
	})

	// PATCH /tasks/:id/toggle -- toggle done status
	kruda.Patch[GetTaskInput, TaskResponse](app, "/tasks/:id/toggle", func(c *kruda.C[GetTaskInput]) (*TaskResponse, error) {
		task, ok := store.ToggleDone(c.In.ID)
		if !ok {
			return nil, kruda.NotFound("task not found")
		}
		return &task, nil
	})

	app.Compile()
	return app
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	store := NewTaskStore()
	app := setupAppWithStore(store)

	log.Println("Task API starting on :3000 ...")
	log.Println("  GET    /tasks           -- list all tasks")
	log.Println("  POST   /tasks           -- create a task")
	log.Println("  GET    /tasks/:id       -- get task by ID")
	log.Println("  PATCH  /tasks/:id/toggle -- toggle done status")
	log.Fatal(app.Listen(":3000"))
}
