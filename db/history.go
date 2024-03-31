package db

import (
	"banking/db/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type History struct {
	dbPool *pgxpool.Pool
}

func NewHistory(dbPool *pgxpool.Pool) *History {
	return &History{
		dbPool: dbPool,
	}
}

func (h *History) GetHistory(ctx context.Context, userId string, limit, offset int) ([]entity.HistoryResponse, error) {
	conn, err := h.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	sql := `
	SELECT id, balance, currency, transfer_proof_image, created_at, source 
	FROM histories 
	WHERE user_id = $1 
	ORDER BY created_at DESC 
	LIMIT $2 OFFSET $3
	`

	rows, err := conn.Query(ctx, sql, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []entity.HistoryResponse
	for rows.Next() {
		var history entity.HistoryResponse
		err = rows.Scan(&history.TransactionId, &history.Balance, &history.Currency, &history.TransferProofImg, &history.CreatedAt, &history.Source)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}

	defer conn.Release()
	return histories, nil
}

func (h *History) GetHistoryCount(ctx context.Context, userId string) (int, error) {
	conn, err := h.dbPool.Acquire(ctx)
	if err != nil {
		return 0, err
	}

	sql := `SELECT COUNT(id) FROM histories WHERE user_id = $1`
	row := conn.QueryRow(ctx, sql, userId)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	defer conn.Release()
	return count, nil
}
