package pgxutils

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPool(ctx context.Context, connURL string) (*pgxpool.Pool, error) {
	pgxPool, err := pgxpool.New(ctx, connURL)
	if err != nil {
		return nil, err
	}

	err = pingWithTimeout(ctx, pgxPool)
	if err != nil {
		return nil, err
	}

	return pgxPool, nil
}

func ConnectPoolWithConfig(ctx context.Context, conf *pgxpool.Config) (*pgxpool.Pool, error) {
	pgxPool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, err
	}

	conn, _ := pgxPool.Acquire(ctx)
	conn.Conn()
	err = pingWithTimeout(ctx, pgxPool)
	if err != nil {
		return nil, err
	}

	return pgxPool, nil
}
