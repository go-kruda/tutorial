package main

import (
	"log"

	"github.com/go-kruda/kruda"
	"github.com/go-kruda/tutorial/04-advanced/08-architecture/complete/handler"
	"github.com/go-kruda/tutorial/04-advanced/08-architecture/complete/repository"
	"github.com/go-kruda/tutorial/04-advanced/08-architecture/complete/service"
)

// ============================================================
// Clean Architecture with DI Container
// ============================================================
//
// This application demonstrates clean architecture using Kruda's
// built-in DI Container. The code is organised into three layers:
//
//   handler/    -> HTTP route handlers (outermost)
//   service/    -> Business logic and validation
//   repository/ -> Data access (innermost)
//
// The dependency direction is always inward:
//   handler -> service -> repository

func main() {
	// -- 1. Create the DI Container --
	container := kruda.NewContainer()

	// -- 2. Register Repository (no deps) --
	repo := repository.NewUserRepository()
	container.Give(repo)

	// -- 3. Register Service (depends on repo) --
	svc := service.NewUserService(repo)
	container.Give(svc)

	// -- 4. Resolve the Service --
	userService := kruda.MustUse[*service.UserService](container)

	// -- 5. Create the Kruda Application --
	app := kruda.New(kruda.WithContainer(container))

	// -- 6. Register Routes --
	handler.RegisterRoutes(app, userService)

	// -- 7. Start the Server --
	log.Println("Architecture demo starting on :3000 ...")
	log.Fatal(app.Listen(":3000"))
}
