package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Clean Architecture with DI Container -- Starter
// ============================================================
//
// In this tutorial you will refactor a single-file application
// into a clean architecture layout with three layers:
//
//   handler/    -> HTTP route handlers
//   service/    -> Business logic and validation
//   repository/ -> Data access (in-memory store)
//
// For the starter, all code lives in main.go so you can see
// the full picture. The complete/ version splits these into
// separate packages.

// ============================================================
// Types
// ============================================================

// User represents a user entity.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUserInput represents the JSON body for creating a user.
type CreateUserInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// GetUserInput captures the :id path parameter.
type GetUserInput struct {
	ID int `param:"id"`
}

// MessageResponse is a generic message envelope.
type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// Repository Layer -- Data Access
// ============================================================

// UserRepository handles persistence operations for users.
type UserRepository struct {
	mu     sync.Mutex
	users  []User
	nextID int
}

// NewUserRepository creates a UserRepository with an empty store.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:  make([]User, 0),
		nextID: 1,
	}
}

// FindAll returns every user in the store.
func (r *UserRepository) FindAll() []User {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]User, len(r.users))
	copy(result, r.users)
	return result
}

// FindByID returns a single user or an error if not found.
func (r *UserRepository) FindByID(id int) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}
	return User{}, fmt.Errorf("user with id %d not found", id)
}

// Create persists a new user and returns the created record.
func (r *UserRepository) Create(name, email string) User {
	r.mu.Lock()
	defer r.mu.Unlock()
	user := User{ID: r.nextID, Name: name, Email: email}
	r.nextID++
	r.users = append(r.users, user)
	return user
}

// Delete removes a user by ID. Returns an error if not found.
func (r *UserRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, u := range r.users {
		if u.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("user with id %d not found", id)
}

// ============================================================
// Service Layer -- Business Logic
// ============================================================

// UserService provides business operations on users.
type UserService struct {
	repo *UserRepository
}

// NewUserService creates a UserService with the given repository.
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

// ListUsers returns all users.
func (s *UserService) ListUsers() []User {
	return s.repo.FindAll()
}

// GetUser returns a single user by ID.
func (s *UserService) GetUser(id int) (User, error) {
	return s.repo.FindByID(id)
}

// CreateUser validates input and creates a new user.
func (s *UserService) CreateUser(name, email string) (User, error) {
	// TODO: Validate that name and email are not empty.
	// Return an error if either is missing.
	//
	// Hint:
	//   if name == "" {
	//       return User{}, fmt.Errorf("name is required")
	//   }
	//   if email == "" {
	//       return User{}, fmt.Errorf("email is required")
	//   }
	return s.repo.Create(name, email), nil
}

// DeleteUser removes a user by ID.
func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}

// ============================================================
// Handler Layer -- Route Registration
// ============================================================

// registerRoutes registers all user-related routes on the app.
func registerRoutes(app *kruda.App, svc *UserService) {
	// TODO: Register GET /users route that returns all users.
	//
	// Hint:
	//   kruda.Get[struct{}, []User](app, "/users", func(c *kruda.C[struct{}]) (*[]User, error) {
	//       users := svc.ListUsers()
	//       return &users, nil
	//   })

	// TODO: Register POST /users route that creates a user.
	// Use CreateUserInput for input binding.
	//
	// Hint:
	//   kruda.Post[CreateUserInput, User](app, "/users", func(c *kruda.C[CreateUserInput]) (*User, error) {
	//       user, err := svc.CreateUser(c.In.Name, c.In.Email)
	//       if err != nil {
	//           return nil, kruda.BadRequest(err.Error())
	//       }
	//       return &user, nil
	//   })

	// TODO: Register GET /users/:id route that returns a user.
	// Use GetUserInput for path parameter binding.
	//
	// Hint:
	//   kruda.Get[GetUserInput, User](app, "/users/:id", func(c *kruda.C[GetUserInput]) (*User, error) {
	//       user, err := svc.GetUser(c.In.ID)
	//       if err != nil {
	//           return nil, kruda.NotFound(err.Error())
	//       }
	//       return &user, nil
	//   })

	// TODO: Register DELETE /users/:id route that deletes a user.
	//
	// Hint:
	//   kruda.Delete[GetUserInput, MessageResponse](app, "/users/:id", func(c *kruda.C[GetUserInput]) (*MessageResponse, error) {
	//       if err := svc.DeleteUser(c.In.ID); err != nil {
	//           return nil, kruda.NotFound(err.Error())
	//       }
	//       return &MessageResponse{
	//           Message: fmt.Sprintf("user %d deleted", c.In.ID),
	//       }, nil
	//   })

	// Keep import used to avoid compile errors.
	_ = svc
}

// ============================================================
// Application Entry Point -- Composition Root
// ============================================================

func main() {
	// TODO: Create a DI Container and register services.
	//
	// Example:
	//   container := kruda.NewContainer()
	//   container.Give(NewUserRepository())
	//   container.GiveLazy(func() (*UserService, error) {
	//       repo := kruda.MustUse[*UserRepository](container)
	//       return NewUserService(repo), nil
	//   })
	var container *kruda.Container
	_ = container

	// TODO: Create a Kruda app with the container and register routes.
	//
	// Example:
	//   app := kruda.New(kruda.WithContainer(container))
	//   registerRoutes(app, kruda.MustUse[*UserService](container))
	//   log.Fatal(app.Listen(":3000"))

	// Keep imports and variables used to avoid compile errors.
	_ = registerRoutes
	_ = NewUserRepository
	_ = NewUserService

	log.Println("Server starting on :3000 ...")
}
