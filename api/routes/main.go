package routes

import (
	"banking/api/handlers"
	"banking/db"

	"github.com/gofiber/fiber/v2"
)

func RouteRegister(app *fiber.App, deps handlers.Dependencies) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	userHandler := handlers.User{
		Database: db.NewUser(deps.DbPool, deps.Cfg),
	}

	UserRoutes(app, userHandler)

	imageUploaderHandler := handlers.ImageUploader{
		Uploader: db.NewImageUploader(deps.Cfg),
	}

	ImageRoutes(app, imageUploaderHandler)
}
