package services

import (
	"context"
	"errors"

	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/repo"
	"github.com/mtchuikov/shortener/pkg/strgen"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
)

type shortenerRepo interface {
	InsertShortenURL(
		context.Context,
		models.ShortenID,
		models.OriginalURL,
	) (
		models.ShortenID,
		error,
	)
}

type shortener struct {
	log     *zerolog.Logger
	baseURL string
	repo    shortenerRepo
}

func NewShortener(log *zerolog.Logger, baseURL string, repo shortenerRepo) *shortener {
	return &shortener{
		log:     log,
		baseURL: baseURL,
		repo:    repo,
	}
}

const shortenerServeOp = "services.shortener.serve"

func (s *shortener) Serve(
	ctx context.Context,
	originalURL models.OriginalURL,
) (
	models.ShortenURL,
	error,
) {
	rawShortenID := strgen.Global.Generate(8)
	shortenID, _ := models.NewShortenID(rawShortenID)

	var err error
	shortenID, err = s.repo.InsertShortenURL(ctx, shortenID, originalURL)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			s.log.Error().Err(err).Str("op", shortenerServeOp).
				Msg("failed to insert shorten url")

			return "", err
		}

		if err == repo.ErrShortenIDAlreadyExists {
			return "", err
		}
	}

	rawShortenURL := s.baseURL + shortenID.String()
	shortenURL, _ := models.NewShortenURL(rawShortenURL)

	return shortenURL, err
}
