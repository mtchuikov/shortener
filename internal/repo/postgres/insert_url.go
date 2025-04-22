package postgres

import (
	"context"
	"fmt"

	"github.com/mtchuikov/shortener/internal/gen/sqlc/v1shortenurls/v1"
	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/repo"
	"github.com/mtchuikov/shortener/pkg/pgcodes"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *shortenURLs) InsertShortenURL(
	ctx context.Context,
	shortenID models.ShortenID,
	originalURL models.OriginalURL,
) (
	models.ShortenID,
	error,
) {
	params := v1shortenurls.InsertShortenURLParams{
		ShortenID:   shortenID.String(),
		OriginalUrl: originalURL.String(),
	}

	rawShortenID, err := r.querier.InsertShortenURL(ctx, params)
	if err != nil {
		pgxErr, ok := err.(*pgconn.PgError)
		if !ok {
			return "", fmt.Errorf("%w: %w", repo.ErrUnexpectedError, err)
		}

		if pgxErr.Code == pgcodes.UniqueViolation {
			return "", repo.ErrShortenIDAlreadyExists
		}

		return "", fmt.Errorf("%w: %w", repo.ErrUnexpectedError, err)
	}

	newShortenID, _ := models.NewShortenID(rawShortenID)
	if newShortenID != shortenID {
		return newShortenID, repo.ErrOirignalURLAlreadyShortened
	}

	return newShortenID, nil
}

func (r *shortenURLs) BatchInsertShortenURLs(
	ctx context.Context, shortURLs models.BatchShortURLs,
) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", repo.ErrUnexpectedError, err)
	}
	defer tx.Rollback(ctx)

	querier := r.querier.WithTx(tx)
	for _, shortURL := range shortURLs {
		params := v1shortenurls.InsertShortenURLParams{
			ShortenID:   shortURL.CorrelationID,
			OriginalUrl: shortURL.OriginalURL,
		}

		_, err = querier.InsertShortenURL(ctx, params)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", repo.ErrUnexpectedError, err)
	}

	return nil
}
