package main

import (
	"testing"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Unit Tests using TestClient
// ============================================================
//
// This file demonstrates how to test a Kruda application using
// Go's standard testing package and kruda.NewTestClient().
//
// Key testing patterns shown:
//   1. Table-driven tests -- the idiomatic Go way to test
//      multiple scenarios with minimal code duplication.
//   2. kruda.NewTestClient(app) -- creates a test client so you
//      can exercise the full HTTP pipeline (routing, input
//      binding, handler logic, JSON serialisation) without
//      starting a real HTTP server.
//   3. Fresh app per test -- each test builds its own app via
//      setupApp() for complete isolation.
//   4. Testing both success and error cases -- good tests
//      verify that errors return the correct HTTP status codes.

// ----------------------------------------------------------------
// Create Task Tests
// ----------------------------------------------------------------

// TestCreateTask_Success verifies that creating a task with valid
// input returns 200 and the correct TaskResponse.
func TestCreateTask_Success(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	resp, err := client.Post("/tasks", map[string]string{
		"title":       "Write unit tests",
		"description": "Learn how to test with TestClient",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if resp.StatusCode() != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode())
	}

	var task TaskResponse
	resp.JSON(&task)

	if task.ID != 1 {
		t.Errorf("expected ID=1, got ID=%d", task.ID)
	}
	if task.Title != "Write unit tests" {
		t.Errorf("expected Title=%q, got %q", "Write unit tests", task.Title)
	}
	if task.Description != "Learn how to test with TestClient" {
		t.Errorf("expected Description=%q, got %q", "Learn how to test with TestClient", task.Description)
	}
	if task.Done != false {
		t.Errorf("expected Done=false, got Done=%v", task.Done)
	}
}

// TestCreateTask_EmptyTitle verifies that creating a task
// without a title returns 400.
func TestCreateTask_EmptyTitle(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	resp, _ := client.Post("/tasks", map[string]string{
		"title":       "",
		"description": "No title provided",
	})
	if resp.StatusCode() != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode())
	}
}

// TestCreateTask_TableDriven demonstrates the table-driven test
// pattern -- the idiomatic Go approach for testing multiple
// input/output combinations in a single test function.
func TestCreateTask_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		title      string
		desc       string
		wantStatus int
		wantTitle  string
	}{
		{
			name:       "valid task with description",
			title:      "Buy groceries",
			desc:       "Milk, eggs, bread",
			wantStatus: 200,
			wantTitle:  "Buy groceries",
		},
		{
			name:       "valid task without description",
			title:      "Quick note",
			desc:       "",
			wantStatus: 200,
			wantTitle:  "Quick note",
		},
		{
			name:       "empty title returns 400",
			title:      "",
			desc:       "This should fail",
			wantStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fresh app per sub-test for isolation.
			app := setupApp()
			client := kruda.NewTestClient(app)

			resp, err := client.Post("/tasks", map[string]string{
				"title":       tt.title,
				"description": tt.desc,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.StatusCode() != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode(), tt.wantStatus)
			}

			if tt.wantStatus == 200 {
				var task TaskResponse
				resp.JSON(&task)
				if task.Title != tt.wantTitle {
					t.Errorf("Title = %q, want %q", task.Title, tt.wantTitle)
				}
				if task.Done != false {
					t.Errorf("Done = %v, want false", task.Done)
				}
			}
		})
	}
}

// ----------------------------------------------------------------
// List Tasks Tests
// ----------------------------------------------------------------

// TestListTasks_Empty verifies that an empty store returns an
// empty JSON array.
func TestListTasks_Empty(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	resp, err := client.Get("/tasks")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.StatusCode() != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode())
	}

	var tasks []TaskResponse
	resp.JSON(&tasks)
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

// TestListTasks_AfterCreate verifies that tasks appear in the
// list after creation.
func TestListTasks_AfterCreate(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	// Create two tasks via the API.
	client.Post("/tasks", map[string]string{"title": "Task A", "description": "First task"})
	client.Post("/tasks", map[string]string{"title": "Task B", "description": "Second task"})

	// List all tasks.
	resp, err := client.Get("/tasks")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.StatusCode() != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode())
	}

	var tasks []TaskResponse
	resp.JSON(&tasks)
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Title != "Task A" {
		t.Errorf("first task Title = %q, want %q", tasks[0].Title, "Task A")
	}
	if tasks[1].Title != "Task B" {
		t.Errorf("second task Title = %q, want %q", tasks[1].Title, "Task B")
	}
}

// ----------------------------------------------------------------
// Get Task By ID Tests
// ----------------------------------------------------------------

// TestGetTaskByID verifies that we can retrieve a task by its ID.
func TestGetTaskByID(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	// Create a task first.
	client.Post("/tasks", map[string]string{"title": "Existing task"})

	// Get it by ID.
	resp, err := client.Get("/tasks/1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode() != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode())
	}

	var task TaskResponse
	resp.JSON(&task)
	if task.Title != "Existing task" {
		t.Errorf("Title = %q, want %q", task.Title, "Existing task")
	}
}

// TestGetTaskByID_NotFound verifies that requesting a
// non-existent task returns 404.
func TestGetTaskByID_NotFound(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	resp, _ := client.Get("/tasks/999")
	if resp.StatusCode() != 404 {
		t.Errorf("expected status 404, got %d", resp.StatusCode())
	}
}

// ----------------------------------------------------------------
// Toggle Done Tests
// ----------------------------------------------------------------

// TestToggleDone_FlipsTwice verifies that toggling a task twice
// returns it to its original state.
func TestToggleDone_FlipsTwice(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	// Create a task.
	client.Post("/tasks", map[string]string{"title": "Toggle me"})

	// First toggle: false -> true
	resp1, err := client.Patch("/tasks/1/toggle", nil)
	if err != nil {
		t.Fatalf("first toggle error: %v", err)
	}
	if resp1.StatusCode() != 200 {
		t.Errorf("expected status 200, got %d", resp1.StatusCode())
	}
	var task1 TaskResponse
	resp1.JSON(&task1)
	if task1.Done != true {
		t.Errorf("after first toggle: Done = %v, want true", task1.Done)
	}

	// Second toggle: true -> false
	resp2, err := client.Patch("/tasks/1/toggle", nil)
	if err != nil {
		t.Fatalf("second toggle error: %v", err)
	}
	var task2 TaskResponse
	resp2.JSON(&task2)
	if task2.Done != false {
		t.Errorf("after second toggle: Done = %v, want false", task2.Done)
	}
}

// TestToggleDone_NotFound verifies that toggling a non-existent
// task returns 404.
func TestToggleDone_NotFound(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	resp, _ := client.Patch("/tasks/42/toggle", nil)
	if resp.StatusCode() != 404 {
		t.Errorf("expected status 404, got %d", resp.StatusCode())
	}
}

// ----------------------------------------------------------------
// TaskStore Unit Tests
// ----------------------------------------------------------------
//
// Testing the store directly is also valuable -- it verifies
// the data layer independently of the HTTP handler layer.

// TestTaskStore_CreateAssignsIncrementingIDs verifies that
// each new task gets a unique, incrementing ID.
func TestTaskStore_CreateAssignsIncrementingIDs(t *testing.T) {
	store := NewTaskStore()

	t1 := store.Create("First", "")
	t2 := store.Create("Second", "")
	t3 := store.Create("Third", "")

	if t1.ID != 1 || t2.ID != 2 || t3.ID != 3 {
		t.Errorf("IDs = [%d, %d, %d], want [1, 2, 3]", t1.ID, t2.ID, t3.ID)
	}
}

// TestTaskStore_FindByID_NotFound verifies that looking up a
// non-existent ID returns false.
func TestTaskStore_FindByID_NotFound(t *testing.T) {
	store := NewTaskStore()

	_, found := store.FindByID(1)
	if found {
		t.Error("expected found=false for empty store")
	}
}
