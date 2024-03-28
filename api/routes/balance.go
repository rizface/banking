package routes

import (
	"banking/api/handlers"
	"banking/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func BalanceRoutes(app *fiber.App, balanceHandler handlers.BalanceHandler) {
	g := app.Group("/v1/balance").Use(middleware.JWTAuth())
	g.Post("", balanceHandler.AddBalance)
	g.Get("", balanceHandler.GetBalances)
	g.Get("/history", balanceHandler.GetHistory)
}
