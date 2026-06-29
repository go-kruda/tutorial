package main

import (
	"testing"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Unit Tests using TestClient + typed helpers
// ============================================================
//
// This file tests the API with kruda.NewTestClient(app) and the
// typed request helpers added in kruda v1.5.0:
//
//   kruda.PostTyped[In, Out](client, path, body)
//   kruda.GetTyped[Out](client, path)
//   kruda.PatchTyped[In, Out](client, path, body)
//
// Each returns a *kruda.TypedTestResponse[Out] whose .Body field is
// the already-decoded response — no manual resp.JSON(&v) needed.
//
// IMPORTANT: the typed helpers are status-agnostic. They decode any
// non-empty JSON body regardless of HTTP status, so always check
// resp.StatusCode() before trusting resp.Body.

// ----------------------------------------------------------------
// Create Task Tests
// ----------------------------------------------------------------

func TestCreateTask_Success(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	resp, err := kruda.PostTyped[CreateTaskInput, TaskResponse](client, "/tasks", CreateTaskInput{
		Title:       "Write unit tests",
		Description: "Learn how to test with TestClient",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.StatusCode() != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode())
	}

	// resp.Body is already a TaskResponse.
	if resp.Body.ID != 1 {
		t.Errorf("expected ID=1, got ID=%d", resp.Body.ID)
	}
	if resp.Body.Title != "Write unit tests" {
		t.Errorf("expected Title=%q, got %q", "Write unit tests", resp.Body.Title)
	}
	if resp.Body.Description != "Learn how to test with TestClient" {
		t.Errorf("expected Description=%q, got %q", "Learn how to test with TestClient", resp.Body.Description)
	}
	if resp.Body.Done != false {
		t.Errorf("expected Done=false, got Done=%v", resp.Body.Done)
	}
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	// On a 400 the error body ({code,message}) does not decode into
	// TaskResponse, so resp.Body stays the zero value — assert on status.
	resp, err := kruda.PostTyped[CreateTaskInput, TaskResponse](client, "/tasks", CreateTaskInput{
		Title:       "",
		Description: "No title provided",
	})
	if err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}
	if resp.StatusCode() != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode())
	}
}

func TestCreateTask_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		title      string
		desc       string
		wantStatus int
		wantTitle  string
	}{
		{name: "valid task with description", title: "Buy groceries", desc: "Milk, eggs, bread", wantStatus: 200, wantTitle: "Buy groceries"},
		{name: "valid task without description", title: "Quick note", desc: "", wantStatus: 200, wantTitle: "Quick note"},
		{name: "empty title returns 400", title: "", desc: "This should fail", wantStatus: 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupApp()
			client := kruda.NewTestClient(app)

			resp, err := kruda.PostTyped[CreateTaskInput, TaskResponse](client, "/tasks", CreateTaskInput{
				Title:       tt.title,
				Description: tt.desc,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.StatusCode() != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode(), tt.wantStatus)
			}
			if tt.wantStatus == 200 {
				if resp.Body.Title != tt.wantTitle {
					t.Errorf("Title = %q, want %q", resp.Body.Title, tt.wantTitle)
				}
				if resp.Body.Done != false {
					t.Errorf("Done = %v, want false", resp.Body.Done)
				}
			}
		})
	}
}

// ----------------------------------------------------------------
// List Tasks Tests
// ----------------------------------------------------------------

func TestListTasks_Empty(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	resp, err := kruda.GetTyped[[]TaskResponse](client, "/tasks")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.StatusCode() != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode())
	}
	if len(resp.Body) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(resp.Body))
	}
}

func TestListTasks_AfterCreate(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	if _, err := kruda.PostTyped[CreateTaskInput, TaskResponse](client, "/tasks", CreateTaskInput{Title: "Task A", Description: "First task"}); err != nil {
		t.Fatalf("seed A: %v", err)
	}
	if _, err := kruda.PostTyped[CreateTaskInput, TaskResponse](client, "/tasks", CreateTaskInput{Title: "Task B", Description: "Second task"}); err != nil {
		t.Fatalf("seed B: %v", err)
	}

	resp, err := kruda.GetTyped[[]TaskResponse](client, "/tasks")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.StatusCode() != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode())
	}
	if len(resp.Body) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(resp.Body))
	}
	if resp.Body[0].Title != "Task A" {
		t.Errorf("first task Title = %q, want %q", resp.Body[0].Title, "Task A")
	}
	if resp.Body[1].Title != "Task B" {
		t.Errorf("second task Title = %q, want %q", resp.Body[1].Title, "Task B")
	}
}

// ----------------------------------------------------------------
// Get Task By ID Tests
// ----------------------------------------------------------------

func TestGetTaskByID(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	if _, err := kruda.PostTyped[CreateTaskInput, TaskResponse](client, "/tasks", CreateTaskInput{Title: "Existing task"}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	resp, err := kruda.GetTyped[TaskResponse](client, "/tasks/1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode() != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode())
	}
	if resp.Body.Title != "Existing task" {
		t.Errorf("Title = %q, want %q", resp.Body.Title, "Existing task")
	}
}

func TestGetTaskByID_NotFound(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	resp, err := kruda.GetTyped[TaskResponse](client, "/tasks/999")
	if err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}
	if resp.StatusCode() != 404 {
		t.Errorf("expected status 404, got %d", resp.StatusCode())
	}
}

// ----------------------------------------------------------------
// Toggle Done Tests
// ----------------------------------------------------------------

func TestToggleDone_FlipsTwice(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	if _, err := kruda.PostTyped[CreateTaskInput, TaskResponse](client, "/tasks", CreateTaskInput{Title: "Toggle me"}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// /toggle takes no request body; pass an empty struct as the typed body.
	resp1, err := kruda.PatchTyped[struct{}, TaskResponse](client, "/tasks/1/toggle", struct{}{})
	if err != nil {
		t.Fatalf("first toggle error: %v", err)
	}
	if resp1.StatusCode() != 200 {
		t.Errorf("expected status 200, got %d", resp1.StatusCode())
	}
	if resp1.Body.Done != true {
		t.Errorf("after first toggle: Done = %v, want true", resp1.Body.Done)
	}

	resp2, err := kruda.PatchTyped[struct{}, TaskResponse](client, "/tasks/1/toggle", struct{}{})
	if err != nil {
		t.Fatalf("second toggle error: %v", err)
	}
	if resp2.Body.Done != false {
		t.Errorf("after second toggle: Done = %v, want false", resp2.Body.Done)
	}
}

func TestToggleDone_NotFound(t *testing.T) {
	app := setupApp()
	client := kruda.NewTestClient(app)

	resp, err := kruda.PatchTyped[struct{}, TaskResponse](client, "/tasks/42/toggle", struct{}{})
	if err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}
	if resp.StatusCode() != 404 {
		t.Errorf("expected status 404, got %d", resp.StatusCode())
	}
}

// ----------------------------------------------------------------
// TaskStore Unit Tests (data layer, no HTTP) — unchanged
// ----------------------------------------------------------------

func TestTaskStore_CreateAssignsIncrementingIDs(t *testing.T) {
	store := NewTaskStore()

	t1 := store.Create("First", "")
	t2 := store.Create("Second", "")
	t3 := store.Create("Third", "")

	if t1.ID != 1 || t2.ID != 2 || t3.ID != 3 {
		t.Errorf("IDs = [%d, %d, %d], want [1, 2, 3]", t1.ID, t2.ID, t3.ID)
	}
}

func TestTaskStore_FindByID_NotFound(t *testing.T) {
	store := NewTaskStore()

	_, found := store.FindByID(1)
	if found {
		t.Error("expected found=false for empty store")
	}
}
