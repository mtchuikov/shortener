package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/repo"
)

type postgres struct {
	conn *pgxpool.Conn
}

func New(conn *pgxpool.Conn) *postgres {
	return &postgres{conn}
}

const createTableQuery = `
	CREATE TABLE IF NOT EXISTS shorten_urls (
		short_id VARCHAR(255) PRIMARY KEY,
		original_url TEXT NOT NULL
	);
`

func (r *postgres) CreateTable(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, createTableQuery)
	return err
}

const insertShortenURLQuery = `
	INSERT INTO shorten_urls (short_id, original_url)
	VALUES ($1, $2)
	ON CONFLICT (short_id) DO NOTHING
`

func (r *postgres) CreateShortURL(ctx context.Context, originalURL, shortID string) error {
	_, err := r.conn.Exec(ctx, insertShortenURLQuery, shortID, originalURL)
	return err
}

func (r *postgres) BatchCreateShortURLs(ctx context.Context, urlsToShort models.URLsToShort) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, url := range urlsToShort {
		_, err = tx.Exec(ctx, insertShortenURLQuery, url.CorrelationID, url.OriginalURL)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

const getOriginalURLQuery = `
	SELECT original_url 
	FROM shorten_urls 
	WHERE short_id = $1`

func (r *postgres) GetOriginalURL(ctx context.Context, shortID string) (string, error) {
	result := r.conn.QueryRow(ctx, getOriginalURLQuery, shortID)

	var originalURL string
	err := result.Scan(&originalURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", repo.ErrOriginalURLNotFound
		}

		return "", err
	}

	return originalURL, nil
}

const getShortIDQuery = `
	SELECT short_id FROM shorten_urls
	WHERE original_url = $1
`

func (r *postgres) GetShortID(ctx context.Context, originalURL string) (string, error) {
	result := r.conn.QueryRow(ctx, getShortIDQuery, originalURL)

	var shortID string
	err := result.Scan(&shortID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", repo.ErrShortIDNotFound
		}

		return "", err
	}

	return shortID, nil
}
