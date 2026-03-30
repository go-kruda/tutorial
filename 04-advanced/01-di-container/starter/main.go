package main

import (
	"fmt"
	"log"

	"github.com/go-kruda/kruda"
)

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
// UserRepository -- Data Access Layer
// ============================================================

type UserRepository struct {
	users  []UserResponse
	nextID int
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:  make([]UserResponse, 0),
		nextID: 1,
	}
}

func (r *UserRepository) FindAll() []UserResponse {
	result := make([]UserResponse, len(r.users))
	copy(result, r.users)
	return result
}

func (r *UserRepository) FindByID(id int) (UserResponse, error) {
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}
	return UserResponse{}, fmt.Errorf("user with id %d not found", id)
}

func (r *UserRepository) Create(name, email string) UserResponse {
	user := UserResponse{ID: r.nextID, Name: name, Email: email}
	r.nextID++
	r.users = append(r.users, user)
	return user
}

func (r *UserRepository) Delete(id int) error {
	for i, u := range r.users {
		if u.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("user with id %d not found", id)
}

// ============================================================
// UserService -- Business Logic Layer
// ============================================================

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ListUsers() []UserResponse { return s.repo.FindAll() }

func (s *UserService) GetUser(id int) (UserResponse, error) { return s.repo.FindByID(id) }

func (s *UserService) CreateUser(name, email string) (UserResponse, error) {
	if name == "" {
		return UserResponse{}, fmt.Errorf("name is required")
	}
	if email == "" {
		return UserResponse{}, fmt.Errorf("email is required")
	}
	return s.repo.Create(name, email), nil
}

func (s *UserService) DeleteUser(id int) error { return s.repo.Delete(id) }

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	// TODO: Create a DI Container
	//
	// Hint:
	//   container := kruda.NewContainer()
	//   container.Give(instance) -- register a singleton
	//   kruda.MustUse[T](container) -- resolve from the container

	// TODO: Register repository and service in the container
	//
	// Hint:
	//   repo := NewUserRepository()
	//   container.Give(repo)
	//   svc := NewUserService(repo)
	//   container.Give(svc)

	// TODO: Resolve service from the container
	//
	// Hint:
	//   userService := kruda.MustUse[*UserService](container)

	// TODO: Create a Kruda app and register routes
	//
	// Pattern A (simple): resolve at startup and use a closure
	//   app := kruda.New()
	//   userService := kruda.MustUse[*UserService](container)
	//
	// Pattern B (recommended): attach the container and resolve per-request
	//   app := kruda.New(kruda.WithContainer(container))
	//   // In the handler:
	//   svc := kruda.MustResolve[*UserService](c.Ctx)
	//   users := svc.ListUsers()

	log.Println("Server starting on :3000 ...")
	// Keep kruda import used for compilation.
	_ = kruda.NewContainer
}
