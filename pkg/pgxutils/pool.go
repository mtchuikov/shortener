package pgxutils

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPool(ctx context.Context, connURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf(wFailedToPingPool, err)
	}

	return pool, nil
}

func AcquireConn(ctx context.Context, pool *pgxpool.Pool) (*pgxpool.Conn, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf(wFailedToAcquireConn, err)
	}

	return conn, nil
}
