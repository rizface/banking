package handlers

import (
	"banking/db"
	"banking/db/entity"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Transaction struct {
	DB *db.Transaction
}

func (t *Transaction) Transfer(c *fiber.Ctx) error {
	var req struct {
		RecipientBankAccountNumber string `json:"recipientBankAccountNumber"`
		RecipientBankName          string `json:"recipientBankName"`
		FromCurrency               string `json:"fromCurrency"`
		Balances                   int    `json:"balances"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if req.Balances < 1 {
		return c.Status(http.StatusBadRequest).
			SendString("balances must be greater than 0")
	}

	if len(req.FromCurrency) != 3 {
		return c.Status(http.StatusBadRequest).
			SendString("currency is invalid")
	}

	if req.RecipientBankName == "" || req.RecipientBankAccountNumber == "" {
		return c.Status(http.StatusBadRequest).
			SendString("recipient bank name and recipient bank account number is required")
	}

	err := t.DB.Transfer(c.UserContext(), entity.Transaction{
		SenderId:                   c.Locals("user_id").(string),
		RecipientBankAccountNumber: req.RecipientBankAccountNumber,
		RecipientBankName:          req.RecipientBankName,
		FromCurrency:               req.FromCurrency,
		Balances:                   req.Balances,
	})

	if errors.Is(err, db.ErrinsuficientBalance) {
		return c.Status(http.StatusBadRequest).
			SendString(err.Error())
	}

	if errors.Is(err, db.ErrBalanceNotFound) {
		return c.Status(http.StatusUnprocessableEntity).
			SendString(err.Error())
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).
			SendString(err.Error())
	}

	return nil
}
