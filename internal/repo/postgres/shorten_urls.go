package postgres

import (
	"github.com/mtchuikov/shortener/internal/gen/sqlc/v1shortenurls/v1"

	"github.com/jackc/pgx/v5"
)

type shortenURLs struct {
	conn    *pgx.Conn
	querier *v1shortenurls.Queries
}

func NewShortenURLs(conn *pgx.Conn) *shortenURLs {
	return &shortenURLs{
		conn:    conn,
		querier: v1shortenurls.New(conn),
	}
}
