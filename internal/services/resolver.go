package services

import (
	"context"
	"errors"

	"github.com/mtchuikov/shortener/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
)

type resolverRepo interface {
	GetOriginalURL(context.Context, models.ShortenID) (models.OriginalURL, error)
}

type resolver struct {
	log  *zerolog.Logger
	repo resolverRepo
}

func NewResolver(log *zerolog.Logger, repo resolverRepo) *resolver {
	return &resolver{
		log:  log,
		repo: repo,
	}
}

const resolverServeOp = "services.shortener.serve"

func (s *resolver) Serve(ctx context.Context, shortenID models.ShortenID,
) (
	models.OriginalURL,
	error,
) {
	originalURL, err := s.repo.GetOriginalURL(ctx, shortenID)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			s.log.Error().Err(err).
				Str("op", resolverServeOp).
				Msg("failed to resolve shorten id")
		}
	}

	return originalURL, err
}
