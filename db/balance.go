package db

import (
	"banking/db/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Balance struct {
	dbPool *pgxpool.Pool
}

func NewBalance(dbPool *pgxpool.Pool) *Balance {
	return &Balance{
		dbPool: dbPool,
	}
}

func (b *Balance) AddBalance(ctx context.Context, balance entity.Balance, history entity.History) error {
	conn, err := b.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		// Handle error
		return err
	}
	defer tx.Rollback(ctx)

	// Upsert balance based on currency
	sql := `
	INSERT INTO balances (user_id, balance, currency) VALUES ($1, $2, $3) 
	ON CONFLICT (user_id, currency) 
	DO UPDATE SET balance = (SELECT balance FROM balances WHERE user_id = $1 and currency = $3) + $2
	`
	_, err = tx.Exec(ctx, sql, balance.UserId, balance.Balance, balance.Currency)
	if err != nil {
		return err
	}

	// Insert history
	sql = `INSERT INTO histories (user_id, balance, currency, transfer_proof_image, source) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(ctx, sql, history.UserId, history.Balance, history.Currency, history.TransferProofImg, history.Source)
	if err != nil {
		return err
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()
	return nil
}

// get balances by user id
func (b *Balance) GetBalances(ctx context.Context, userId string) ([]entity.BalanceResponse, error) {
	conn, err := b.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	sql := `SELECT balance, currency FROM balances WHERE user_id = $1 order by balance desc`
	rows, err := conn.Query(ctx, sql, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []entity.BalanceResponse
	for rows.Next() {
		var balance entity.BalanceResponse
		err := rows.Scan(&balance.Balance, &balance.Currency)
		if err != nil {
			return nil, err
		}

		balances = append(balances, balance)
	}

	return balances, nil
}
