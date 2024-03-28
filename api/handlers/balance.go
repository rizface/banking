package handlers

import (
	"banking/api/responses"
	"banking/db"
	"banking/db/entity"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
)

type (
	BalanceHandler struct {
		Database        *db.Balance
		HistoryDatabase *db.History
	}
	AddBalanceRequest struct {
		SenderBankAccountNumber string  `json:"senderBankAccountNumber" `
		SenderBankName          string  `json:"senderBankName" `
		AddedBalance            float64 `json:"addedBalance" `
		Currency                string  `json:"currency"`
		TransferProofImg        string  `json:"transferProofImg" `
	}

	QueryGetHistory struct {
		Limit  int `query:"limit"`
		Offset int `query:"offset"`
	}
)

func (abr AddBalanceRequest) Validate() error {
	return validation.ValidateStruct(&abr,
		validation.Field(&abr.SenderBankAccountNumber, validation.Required, validation.Length(5, 30)),
		validation.Field(&abr.SenderBankName, validation.Required, validation.Length(5, 30)),
		validation.Field(&abr.AddedBalance, validation.Required, validation.Min(float64(1.0))),
		validation.Field(&abr.Currency, validation.Required, is.CurrencyCode),
		validation.Field(&abr.TransferProofImg, validation.Required),
	)
}

func (qg QueryGetHistory) Validate() error {
	return validation.ValidateStruct(&qg,
		validation.Field(&qg.Limit, validation.Min(1)),
		validation.Field(&qg.Offset, validation.Min(0)),
	)
}

func (b *BalanceHandler) AddBalance(ctx *fiber.Ctx) error {
	var req AddBalanceRequest
	if err := ctx.BodyParser(&req); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	// convert request to entity.Balance and entity.History
	balance := entity.Balance{
		UserId:   ctx.Locals("user_id").(string),
		Balance:  req.AddedBalance,
		Currency: req.Currency,
	}

	history := entity.History{
		UserId:           ctx.Locals("user_id").(string),
		Balance:          req.AddedBalance,
		Currency:         req.Currency,
		TransferProofImg: req.TransferProofImg,
		Source: entity.Source{
			BankAccountNumber: req.SenderBankAccountNumber,
			BankName:          req.SenderBankName,
		},
	}

	if err := b.Database.AddBalance(ctx.Context(), balance, history); err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	return responses.ReturnTheResponse(ctx, false, 200, "Balance added successfully", nil)
}

// get balances by user id
func (b *BalanceHandler) GetBalances(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(string)
	balances, err := b.Database.GetBalances(ctx.Context(), userId)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	return responses.ReturnTheResponse(ctx, false, 200, "success", balances)
}

// GetHistory by user id
func (b *BalanceHandler) GetHistory(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(string)

	var qg QueryGetHistory
	if err := ctx.QueryParser(&qg); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	if err := qg.Validate(); err != nil {
		return responses.ErrorBadRequest(ctx, err.Error())
	}

	if qg.Limit == 0 {
		qg.Limit = 5
	}

	histories, err := b.HistoryDatabase.GetHistory(ctx.Context(), userId, qg.Limit, qg.Offset)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	count, err := b.HistoryDatabase.GetHistoryCount(ctx.Context(), userId)
	if err != nil {
		return responses.ErrorInternalServerError(ctx, err.Error())
	}

	return responses.ReturnTheResponseMeta(ctx, false, 200, "success", histories, responses.Meta{
		Total:  count,
		Limit:  qg.Limit,
		Offset: qg.Offset,
	})
}
