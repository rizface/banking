package routes

import (
	"banking/api/handlers"
	"banking/db"
	"banking/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func RouteRegister(app *fiber.App, deps handlers.Dependencies) {
	userHandler := handlers.User{
		Database: db.NewUser(deps.DbPool, deps.Cfg),
	}

	UserRoutes(app, userHandler)

	imageUploaderHandler := handlers.ImageUploader{
		Uploader: utils.NewImageUploader(deps.Cfg),
	}
	ImageRoutes(app, imageUploaderHandler)

	balanceHandler := handlers.BalanceHandler{
		Database:        db.NewBalance(deps.DbPool),
		HistoryDatabase: db.NewHistory(deps.DbPool),
	}

	BalanceRoutes(app, balanceHandler)

	transactionHandler := handlers.Transaction{
		DB: db.NewTransaction(deps.DbPool),
	}

	TransactionRoutes(app, transactionHandler)
}
