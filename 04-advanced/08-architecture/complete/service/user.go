package service

import (
	"fmt"

	"github.com/go-kruda/tutorial/04-advanced/08-architecture/complete/repository"
)

// ============================================================
// 🧩 Service Layer — Business Logic
// ============================================================
//
// The service layer sits between handlers and repositories. It
// contains business rules, validation, and orchestration logic.
// It depends ONLY on the repository layer — never on HTTP
// concepts like requests, responses, or contexts.
//
// Why a separate service package?
// --------------------------------
// Keeping business logic in its own package ensures that:
//   - Handlers stay thin (just HTTP adapters)
//   - Business rules are reusable across different transports
//     (HTTP, gRPC, CLI, etc.)
//   - Testing business logic doesn't require HTTP infrastructure
//   - The dependency direction is always inward:
//     handler → service → repository

// UserResponse is the service-level representation of a user
// returned to callers. In a larger app this might differ from
// the repository entity (e.g., hiding internal fields).
type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserService provides business operations on users.
// It depends on *repository.UserRepository — the DI Container
// injects this dependency at resolution time.
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates a UserService with the given repository.
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// ListUsers returns all users.
func (s *UserService) ListUsers() []UserResponse {
	users := s.repo.FindAll()
	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = UserResponse{ID: u.ID, Name: u.Name, Email: u.Email}
	}
	return result
}

// GetUser returns a single user by ID.
func (s *UserService) GetUser(id int) (UserResponse, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return UserResponse{}, err
	}
	return UserResponse{ID: u.ID, Name: u.Name, Email: u.Email}, nil
}

// CreateUser validates input and creates a new user.
// Business rules: name and email must not be empty.
func (s *UserService) CreateUser(name, email string) (UserResponse, error) {
	if name == "" {
		return UserResponse{}, fmt.Errorf("name is required")
	}
	if email == "" {
		return UserResponse{}, fmt.Errorf("email is required")
	}
	u := s.repo.Create(name, email)
	return UserResponse{ID: u.ID, Name: u.Name, Email: u.Email}, nil
}

// DeleteUser removes a user by ID.
func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}
