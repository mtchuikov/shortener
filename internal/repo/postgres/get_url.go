package postgres

import (
	"context"
	"fmt"

	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/repo"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *shortenURLs) GetOriginalURL(
	ctx context.Context,
	shortenID models.ShortenID,
) (
	models.OriginalURL,
	error,
) {
	rawShortenID := shortenID.String()
	rawOriginalURL, err := r.querier.GetOriginalURL(ctx, rawShortenID)
	if err != nil {
		_, ok := err.(*pgconn.PgError)
		if ok || err != pgx.ErrNoRows {
			return "", fmt.Errorf("%w: %w", repo.ErrUnexpectedError, err)
		}

		return "", repo.ErrOriginalURLNotFound
	}

	originalURL, _ := models.NewOriginalURL(rawOriginalURL)
	return originalURL, nil
}
