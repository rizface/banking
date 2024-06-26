package db

import (
	"banking/configs"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgConn(config configs.Config) (*pgxpool.Pool, error) {
	ctx := context.Background()

	var dsn string

	if os.Getenv("ENV") == "production" {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", config.DbUsername, config.DbPassword, config.DbHost, config.DbPort, config.DbName, config.DbParams)
	} else {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DbUsername, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	}

	dbconfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	dbconfig.MaxConnLifetime = 1 * time.Hour
	dbconfig.MaxConnIdleTime = 30 * time.Minute
	dbconfig.HealthCheckPeriod = 5 * time.Second
	dbconfig.MaxConns = 15
	dbconfig.MinConns = 5

	return pgxpool.NewWithConfig(ctx, dbconfig)
}
