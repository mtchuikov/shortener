package services

import (
	"context"
	"errors"

	"github.com/rs/zerolog"

	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/storage"
	"github.com/mtchuikov/shortener/pkg/strgen"
)

var _ Shortener = (*shortener)(nil)

type shortener struct {
	log     *zerolog.Logger
	storage storage.ShortURLsStorage
	slugGen *strgen.Generator
}

func NewShortener(log *zerolog.Logger, storage storage.ShortURLsStorage) *shortener {
	return &shortener{
		log:     log,
		storage: storage,
		slugGen: strgen.Global,
	}
}

const createShortURLOp = "services.shortener.create_short_url"

func (s *shortener) CreateShortURL(
	ctx context.Context,
	userID,
	originalURL string,
) (string, error) {
	slug := s.slugGen.Generate(8)
	shortURL, err := s.storage.CreateShortURL(ctx, userID, slug, originalURL)
	if err != nil {
		if !errors.Is(err, storage.ErrOriginalURLAlreadyExists) {
			s.log.Error().Err(err).
				Str("op", createShortURLOp).
				Msg("failed to create short url")

			return "", err
		}
	}

	return shortURL, err
}

const batchCreateShortURLsOp = "services.shortener.batch_create_short_urls"

func (s *shortener) BatchCreateShortURLs(
	ctx context.Context,
	userID string,
	items []models.BatchCreateShortURLs,
) ([]models.BatchCreateShortURLsResult, error) {
	result, err := s.storage.BatchCreateShortURLs(ctx, userID, items)
	if err != nil {
		if !errors.Is(err, storage.ErrNoItemsInBatch) {
			s.log.Error().Err(err).
				Str("op", batchCreateShortURLsOp).
				Msg("failed to batch create short urls")
		}

		return nil, err
	}

	return result, nil
}

const getOriginalURLBySlugOp = "services.shortener.get_original_url_by_slug"

func (s *shortener) GetOriginalURLBySlug(
	ctx context.Context,
	slug string,
) (string, bool, error) {
	originalURL, deleted, err := s.storage.GetOriginalURLBySlug(ctx, slug)
	if err != nil {
		if !errors.Is(err, storage.ErrOriginalURLNotFound) {
			s.log.Error().Err(err).
				Str("op", getOriginalURLBySlugOp).
				Msg("failed to get original url")
		}

		return "", deleted, err
	}

	return originalURL, deleted, nil
}

const listShortURLsByUserOp = "services.shortener.list_short_urls_by_user"

func (s *shortener) ListShortURLsByUser(
	ctx context.Context,
	userID string,
) ([]models.ListShortURLsByUserResult, error) {
	result, err := s.storage.ListShortURLsByUser(ctx, userID)
	if err != nil {
		if !errors.Is(err, storage.ErrUserHasNoShortURLs) {
			s.log.Error().Err(err).
				Str("op", listShortURLsByUserOp).
				Msg("failed to list short urls")
		}

		return nil, err
	}

	return result, nil
}

const markShortURLAsActiveOp = "services.shortener.mark_short_url_as_active"

func (s *shortener) MarkShortURLsAsActive(ctx context.Context, userID string, slugs []string) error {
	err := s.storage.MarkShortURLsAsActive(ctx, userID, slugs)
	if err != nil {
		s.log.Error().Err(err).
			Str("op", markShortURLAsActiveOp).
			Msg("failed to set short url status")

		return err
	}

	return nil
}

const markShortURLAsDeletedOp = "services.shortener.mark_short_url_as_deleted"

func (s *shortener) MarkShortURLsAsDeleted(ctx context.Context, userID string, slugs []string) error {
	err := s.storage.MarkShortURLsAsDeleted(ctx, userID, slugs)
	if err != nil {
		s.log.Error().Err(err).
			Str("op", markShortURLAsDeletedOp).
			Msg("failed to set short url status")

		return err
	}

	return nil
}
