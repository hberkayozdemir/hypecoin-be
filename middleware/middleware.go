package middleware

import (
	"github.com/hberkayozdemir/hypecoin-be/auth"
	"github.com/hberkayozdemir/hypecoin-be/repository"
	"github.com/gofiber/fiber/v2"
)

func SetupMiddleWare(app *fiber.App, userRepository repository.Repository) {
	authService := auth.NewService(userRepository)
	authHandler := auth.NewHandler(authService)
	app.Use("/user", authHandler.AuthUserHandler)
	app.Use("/admin", authHandler.AuthAdminHandler)
	app.Use("/editor",authHandler.AuthEditorHandler)
}
