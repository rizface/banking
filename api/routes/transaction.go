package routes

import (
	"banking/api/handlers"
	"banking/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func TransactionRoutes(app *fiber.App, h handlers.Transaction) {
	g := app.Group("/v1/transaction").Use(middleware.JWTAuth())

	g.Post("", h.Transfer)
}
