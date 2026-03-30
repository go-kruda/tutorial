package main

import (
	"log"
	"sync"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Request / Response Types
// ============================================================

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
	// TODO: Lock the mutex, create a TaskResponse with the next
	// ID, append it to the tasks slice, increment nextID, and
	// return the new task.
	//
	// Hint:
	//   s.mu.Lock()
	//   defer s.mu.Unlock()
	//   task := TaskResponse{
	//       ID:          s.nextID,
	//       Title:       title,
	//       Description: description,
	//       Done:        false,
	//   }
	//   s.nextID++
	//   s.tasks = append(s.tasks, task)
	//   return task
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = title
	_ = description
	return TaskResponse{}
}

// FindByID returns a task by ID, or false if not found.
func (s *TaskStore) FindByID(id int) (TaskResponse, bool) {
	// TODO: Lock the mutex, iterate over tasks, and return the
	// matching task.
	//
	// Hint:
	//   s.mu.Lock()
	//   defer s.mu.Unlock()
	//   for _, t := range s.tasks {
	//       if t.ID == id { return t, true }
	//   }
	//   return TaskResponse{}, false
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = id
	return TaskResponse{}, false
}

// ToggleDone flips the Done status of a task by ID.
func (s *TaskStore) ToggleDone(id int) (TaskResponse, bool) {
	// TODO: Lock the mutex, find the task by ID, flip its Done
	// field, and return the updated task.
	//
	// Hint:
	//   s.mu.Lock()
	//   defer s.mu.Unlock()
	//   for i, t := range s.tasks {
	//       if t.ID == id {
	//           s.tasks[i].Done = !t.Done
	//           return s.tasks[i], true
	//       }
	//   }
	//   return TaskResponse{}, false
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = id
	return TaskResponse{}, false
}

// ============================================================
// Route Setup
// ============================================================

// setupApp creates and configures a Kruda app with all routes.
// This function is used in both main() and the test file.
func setupApp() *kruda.App {
	store := NewTaskStore()
	return setupAppWithStore(store)
}

// setupAppWithStore creates a Kruda app using the provided store.
func setupAppWithStore(store *TaskStore) *kruda.App {
	app := kruda.New()

	// TODO: Register GET /tasks route that returns all tasks.
	//
	// Hint:
	//   kruda.Get[struct{}, []TaskResponse](app, "/tasks", func(c *kruda.C[struct{}]) (*[]TaskResponse, error) {
	//       tasks := store.All()
	//       return &tasks, nil
	//   })

	// TODO: Register POST /tasks route that creates a new task.
	// Validate that the title is not empty -- return kruda.BadRequest
	// if it is.
	//
	// Hint:
	//   kruda.Post[CreateTaskInput, TaskResponse](app, "/tasks", func(c *kruda.C[CreateTaskInput]) (*TaskResponse, error) {
	//       if c.In.Title == "" {
	//           return nil, kruda.BadRequest("title is required")
	//       }
	//       task := store.Create(c.In.Title, c.In.Description)
	//       return &task, nil
	//   })

	// TODO: Register GET /tasks/:id route that returns a task by ID.
	// Use GetTaskInput to capture the path parameter.
	//
	// Hint:
	//   kruda.Get[GetTaskInput, TaskResponse](app, "/tasks/:id", func(c *kruda.C[GetTaskInput]) (*TaskResponse, error) {
	//       task, ok := store.FindByID(c.In.ID)
	//       if !ok {
	//           return nil, kruda.NotFound("task not found")
	//       }
	//       return &task, nil
	//   })

	// TODO: Register PATCH /tasks/:id/toggle route that toggles
	// a task's done status.
	//
	// Hint:
	//   kruda.Patch[GetTaskInput, TaskResponse](app, "/tasks/:id/toggle", func(c *kruda.C[GetTaskInput]) (*TaskResponse, error) {
	//       task, ok := store.ToggleDone(c.In.ID)
	//       if !ok {
	//           return nil, kruda.NotFound("task not found")
	//       }
	//       return &task, nil
	//   })

	// Keep imports used to avoid compile errors.
	_ = store

	app.Compile()
	return app
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	// TODO: Build the app and start the server.
	//
	// Example:
	//   app := setupApp()
	//   log.Fatal(app.Listen(":3000"))

	log.Println("Task API starting on :3000 ...")
}
