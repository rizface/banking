package routes

import (
	"banking/api/handlers"
	"banking/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func ImageRoutes(app *fiber.App, h handlers.ImageUploader) {
	app.Post("/v1/image", h.Upload).Use(middleware.JWTAuth())
}
