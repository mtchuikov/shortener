package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	conn *pgxpool.Conn
}

func New(conn *pgxpool.Conn) *repo {
	return &repo{conn: conn}
}

const createTableQuery = `
	CREATE TABLE IF NOT EXISTS shorten_urls (
		id VARCHAR(255) PRIMARY KEY,
		url TEXT NOT NULL UNIQUE
	);
`

func (r *repo) CreateTable(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, createTableQuery)
	return err
}

const insertShortenURLQuery = `
	INSERT INTO shorten_urls (id, url)
	VALUES ($1, $2)
	ON CONFLICT (id) DO NOTHING
`

func (r *repo) CreateShortURL(ctx context.Context, url, id string) error {
	_, err := r.conn.Exec(ctx, insertShortenURLQuery, id, url)
	return err
}

const getOriginalURLQuery = "SELECT url FROM shorten_urls WHERE id = $1"

func (r *repo) GetOriginalURL(ctx context.Context, shortID string) (string, error) {
	result := r.conn.QueryRow(ctx, getOriginalURLQuery, shortID)

	var originalURL string
	err := result.Scan(&originalURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}

		return "", err
	}

	return originalURL, nil
}

const getShortIDQuery = `
	SELECT id FROM shorten_urls
	WHERE url = $1;
`

func (r *repo) GetShortID(ctx context.Context, originalURL string) (string, error) {
	result := r.conn.QueryRow(ctx, getShortIDQuery, originalURL)

	var shortID string
	err := result.Scan(&shortID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}

		return "", err
	}

	return shortID, nil
}
