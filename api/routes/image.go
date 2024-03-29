package routes

import (
	"banking/api/handlers"
	"banking/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func ImageRoutes(app *fiber.App, h handlers.ImageUploader) {
	g := app.Group("/v1/image").Use(middleware.JWTAuth())
	g.Post("", h.Upload)
}
