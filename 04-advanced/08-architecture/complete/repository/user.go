package repository

import (
	"fmt"
	"sync"
)

// ============================================================
// 🗄️ Repository Layer — Data Access
// ============================================================
//
// The repository layer is the lowest layer in clean architecture.
// It encapsulates ALL data access logic — in this tutorial we
// use an in-memory store, but in production you would swap this
// for a database-backed implementation without changing the
// service or handler layers.
//
// Why a separate repository package?
// -----------------------------------
// By isolating data access into its own package, the service
// layer depends only on the repository's exported API. This
// makes it trivial to:
//   - Swap in-memory storage for PostgreSQL, Redis, etc.
//   - Write unit tests with a fresh store per test
//   - Enforce that business logic never touches raw storage

// User represents a user entity in the data layer.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserRepository handles persistence operations for users.
// It uses an in-memory slice protected by a mutex for
// thread-safe concurrent access.
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

	user := User{
		ID:    r.nextID,
		Name:  name,
		Email: email,
	}
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
