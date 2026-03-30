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
	user := UserResponse{
		ID:    r.nextID,
		Name:  name,
		Email: email,
	}
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

func (s *UserService) ListUsers() []UserResponse {
	return s.repo.FindAll()
}

func (s *UserService) GetUser(id int) (UserResponse, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) CreateUser(name, email string) (UserResponse, error) {
	if name == "" {
		return UserResponse{}, fmt.Errorf("name is required")
	}
	if email == "" {
		return UserResponse{}, fmt.Errorf("email is required")
	}
	return s.repo.Create(name, email), nil
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	// ── 1. Create the DI Container ───────────────────────────
	//
	// kruda.NewContainer() creates an empty dependency injection
	// container. You register instances with container.Give()
	// and resolve them with kruda.MustUse[T](container).
	//
	// Why use a DI Container?
	// -----------------------
	// As your application grows, manually wiring dependencies in
	// main() becomes tedious and error-prone. The DI Container
	// centralises construction logic and makes it trivial to swap
	// implementations (e.g., a mock repository for testing).
	container := kruda.NewContainer()

	// ── 2. Register Services ─────────────────────────────────
	//
	// container.Give() registers a singleton by its concrete type.
	// The container stores the instance and resolves it by type
	// when requested.
	repo := NewUserRepository()
	container.Give(repo)

	// Build the service with its dependency (repository).
	svc := NewUserService(repo)
	container.Give(svc)

	// ── 3. Resolve Services ──────────────────────────────────
	//
	// Pattern A: Resolve at startup with MustUse[T].
	// Good for simple apps where you wire everything in main().
	userService := kruda.MustUse[*UserService](container)
	_ = userService // used in Pattern A handlers below

	// ── 4. Create the Kruda Application ──────────────────────
	//
	// Pattern B (recommended): Attach the container to the App
	// with WithContainer, then resolve per-request with
	// MustResolve[T](c). This is the idiomatic Kruda DI pattern
	// — services are resolved from the request context, making
	// handlers testable and decoupled from main().
	app := kruda.New(kruda.WithContainer(container))

	// ── 5. Register Routes ───────────────────────────────────
	//
	// These handlers use Pattern B: kruda.MustResolve[T](c)
	// resolves the service from the request context. The
	// container was attached via WithContainer above.
	kruda.Get[struct{}, []UserResponse](app, "/users", func(c *kruda.C[struct{}]) (*[]UserResponse, error) {
		svc := kruda.MustResolve[*UserService](c.Ctx)
		users := svc.ListUsers()
		return &users, nil
	})

	kruda.Post[CreateUserInput, UserResponse](app, "/users", func(c *kruda.C[CreateUserInput]) (*UserResponse, error) {
		svc := kruda.MustResolve[*UserService](c.Ctx)
		user, err := svc.CreateUser(c.In.Name, c.In.Email)
		if err != nil {
			return nil, kruda.BadRequest(err.Error())
		}
		return &user, nil
	})

	kruda.Get[GetUserInput, UserResponse](app, "/users/:id", func(c *kruda.C[GetUserInput]) (*UserResponse, error) {
		svc := kruda.MustResolve[*UserService](c.Ctx)
		user, err := svc.GetUser(c.In.ID)
		if err != nil {
			return nil, kruda.NotFound(err.Error())
		}
		return &user, nil
	})

	kruda.Delete[GetUserInput, MessageResponse](app, "/users/:id", func(c *kruda.C[GetUserInput]) (*MessageResponse, error) {
		svc := kruda.MustResolve[*UserService](c.Ctx)
		if err := svc.DeleteUser(c.In.ID); err != nil {
			return nil, kruda.NotFound(err.Error())
		}
		return &MessageResponse{
			Message: fmt.Sprintf("user %d deleted", c.In.ID),
		}, nil
	})

	// ── 6. Start the Server ──────────────────────────────────
	log.Println("Server starting on :3000 ...")
	log.Fatal(app.Listen(":3000"))
}
