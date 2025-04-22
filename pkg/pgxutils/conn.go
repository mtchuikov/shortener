package pgxutils

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Conn struct {
	*pgx.Conn
}

func Connect(ctx context.Context, connURL string) (*Conn, error) {
	pgxConn, err := pgx.Connect(ctx, connURL)
	if err != nil {
		return nil, err
	}

	err = pingWithTimeout(ctx, pgxConn)
	if err != nil {
		return nil, err
	}

	return &Conn{Conn: pgxConn}, nil
}

func ConnectWithConfig(ctx context.Context, conf *pgx.ConnConfig) (*Conn, error) {
	pgxConn, err := pgx.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, err
	}

	err = pingWithTimeout(ctx, pgxConn)
	if err != nil {
		return nil, err
	}

	return &Conn{Conn: pgxConn}, nil
}
