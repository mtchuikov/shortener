package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/storage"
)

func (s *shortURLs) GetOriginalURLBySlug(ctx context.Context, slug string) (string, bool, error) {
	result, err := s.queries.GetOriginalURLBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, storage.ErrOriginalURLNotFound
		}

		err = fmt.Errorf("%s: %w", storage.ErrMsgSomethingWentWrong, err)
		return "", false, err
	}

	return result.OriginalUrl, result.Deleted, nil
}

func (s *shortURLs) ListShortURLsByUser(
	ctx context.Context,
	userID string,
) ([]models.ListShortURLsByUserResult, error) {
	urlsOut, err := s.queries.ListShortURLsByUser(ctx, userID)
	if err != nil {
		err = fmt.Errorf("%s: %w", storage.ErrMsgSomethingWentWrong, err)
		return nil, err
	}

	urlsNum := len(urlsOut)
	if urlsNum == 0 {
		return nil, storage.ErrUserHasNoShortURLs
	}

	result := make([]models.ListShortURLsByUserResult, urlsNum)

	for idx, item := range urlsOut {
		result[idx].ShortURL = s.baseURL + item.Slug
		result[idx].OriginalURL = item.OriginalUrl
	}

	return result, nil
}
