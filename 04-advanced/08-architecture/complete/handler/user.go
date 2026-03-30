package handler

import (
	"fmt"

	"github.com/go-kruda/kruda"
	"github.com/go-kruda/tutorial/04-advanced/08-architecture/complete/service"
)

// ============================================================
// Handler Layer -- HTTP Adapters
// ============================================================
//
// The handler layer is the outermost layer in clean architecture.
// Handlers are thin adapters that translate between HTTP and the
// service layer. They:
//   1. Define input structs with json/param tags for automatic binding
//   2. Call the appropriate service method
//   3. Return a typed response (Kruda serialises it to JSON)
//
// Handlers depend ONLY on the service layer -- they have no
// knowledge of repositories or how data is stored. This keeps
// them focused on HTTP concerns only.

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

// RegisterRoutes registers all user-related routes on the app.
// The resolved service is injected -- handlers have no knowledge
// of the repository or the DI container.
func RegisterRoutes(app *kruda.App, svc *service.UserService) {
	// GET /users -- list all users
	kruda.Get[struct{}, []service.UserResponse](app, "/users", func(c *kruda.C[struct{}]) (*[]service.UserResponse, error) {
		users := svc.ListUsers()
		return &users, nil
	})

	// POST /users -- create a new user
	kruda.Post[CreateUserInput, service.UserResponse](app, "/users", func(c *kruda.C[CreateUserInput]) (*service.UserResponse, error) {
		user, err := svc.CreateUser(c.In.Name, c.In.Email)
		if err != nil {
			return nil, kruda.BadRequest(err.Error())
		}
		return &user, nil
	})

	// GET /users/:id -- get a single user
	kruda.Get[GetUserInput, service.UserResponse](app, "/users/:id", func(c *kruda.C[GetUserInput]) (*service.UserResponse, error) {
		user, err := svc.GetUser(c.In.ID)
		if err != nil {
			return nil, kruda.NotFound(err.Error())
		}
		return &user, nil
	})

	// DELETE /users/:id -- delete a user
	kruda.Delete[GetUserInput, MessageResponse](app, "/users/:id", func(c *kruda.C[GetUserInput]) (*MessageResponse, error) {
		if err := svc.DeleteUser(c.In.ID); err != nil {
			return nil, kruda.NotFound(err.Error())
		}
		return &MessageResponse{
			Message: fmt.Sprintf("user %d deleted", c.In.ID),
		}, nil
	})
}
