package db

import (
	"banking/db/entity"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction struct {
	dbPool *pgxpool.Pool
}

func NewTransaction(db *pgxpool.Pool) *Transaction {
	return &Transaction{
		dbPool: db,
	}
}

func (t *Transaction) Transfer(ctx context.Context, trx entity.Transaction) error {
	conn, err := t.dbPool.Acquire(ctx)
	if err != nil {
		return errors.New("failed acquire connection from pool")
	}

	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed start transaction: %v", err)
	}

	var Balance struct {
		Id      int
		Balance int
	}

	err = tx.
		QueryRow(ctx, `select id, balance from balances where user_id = $1 and currency = $2 for update`, trx.SenderId, trx.FromCurrency).
		Scan(
			&Balance.Id, &Balance.Balance,
		)
	if errors.Is(err, pgx.ErrNoRows) {
		tx.Rollback(ctx)
		return ErrBalanceNotFound
	}

	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if Balance.Balance < trx.Balances {
		tx.Rollback(ctx)
		return ErrinsuficientBalance
	}

	_, err = tx.Exec(ctx, `update balances set balance = balance - $1 where id = $2`, trx.Balances, Balance.Id)
	if err != nil {
		return err
	}

	source := entity.Source{
		BankAccountNumber: trx.RecipientBankAccountNumber,
		BankName:          trx.RecipientBankName,
	}

	sql := `
		insert into histories (user_id, balance, currency, transfer_proof_image, source) values (
			$1, $2, $3, $4, $5
		)
	`

	_, err = tx.Exec(ctx, sql, trx.SenderId, -trx.Balances, trx.FromCurrency, "", source)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("failed insert transaction history: %v", err)
	}

	tx.Commit(ctx)

	return nil
}
