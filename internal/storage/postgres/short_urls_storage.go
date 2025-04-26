package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mtchuikov/shortener/internal/gen/sqlc/shorturls"
	"github.com/mtchuikov/shortener/internal/storage"
)

var _ storage.ShortURLsStorage = (*shortURLs)(nil)

type shortURLs struct {
	conn    *pgxpool.Pool
	queries *shorturls.Queries
	baseURL string
}

func NewShortURLs(conn *pgxpool.Pool, baseURL string) *shortURLs {
	return &shortURLs{
		conn:    conn,
		queries: shorturls.New(conn),
		baseURL: baseURL,
	}
}
