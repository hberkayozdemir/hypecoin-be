package main

import (
	"github.com/hberkayozdemir/hypecoin-be/controller"
	"github.com/hberkayozdemir/hypecoin-be/middleware"
	"github.com/hberkayozdemir/hypecoin-be/repository"
	"github.com/hberkayozdemir/hypecoin-be/service"
	"github.com/hberkayozdemir/hypecoin-be/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	dbURL := utils.GetDBUrl()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	repository := repository.NewRepository(dbURL)
	middleware.SetupMiddleWare(app, *repository)
	service := service.NewService(repository)
	api := controller.NewAPI(&service)

	api.SetupApp(app)

	app.Listen(":8080")
}
