package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mtchuikov/shortener/pkg/pinger"
	"github.com/rs/zerolog"
)

type Pool struct {
	log     zerolog.Logger
	connURL string
	pool    *pgxpool.Pool
}

func NewPool(log zerolog.Logger, connURL string) *Pool {
	return &Pool{
		log:     log,
		connURL: connURL,
	}
}

const wPoolPingFailed = "failed to ping postgres: %w"
const wPoolPingTimeout = "failed to ping postgres, timeout: %w"

func ConnectPool(ctx context.Context, connURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = fmt.Errorf(wPoolPingTimeout, err)
			return nil, err
		}

		return nil, fmt.Errorf(wPoolPingFailed, err)
	}

	return pool, nil
}

func (p *Pool) AutoReconnect(ctx context.Context, pinger *pinger.Pinger) {
	errCh := pinger.Subscribe()
	defer pinger.Unsubscribe(errCh)

	for {
		select {
		case <-ctx.Done():
			return
		case <-errCh:
			pool, err := pgxpool.New(ctx, p.connURL)
			if err != nil {
				p.log.Error().
					Err(err).
					Msg("err")
			}

			p.pool = pool
			pinger.ReplacePinger(pool)
		}
	}
}
