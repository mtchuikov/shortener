package postgres

import (
	"context"
	"fmt"

	"github.com/mtchuikov/shortener/internal/gen/sqlc/shorturls"
	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/storage"
)

func (s *shortURLs) CreateShortURL(
	ctx context.Context,
	userID,
	slug,
	originalURL string,
) (string, error) {
	slugOut, err := s.queries.CreateShortURL(ctx,
		shorturls.CreateShortURLParams{
			UserID:      userID,
			Slug:        slug,
			OriginalUrl: originalURL,
		})
	if err != nil {
		err = fmt.Errorf("%s: %w", storage.ErrMsgSomethingWentWrong, err)
		return "", err
	}

	if slugOut != slug {
		err = storage.ErrOriginalURLAlreadyExists
	}

	return s.baseURL + slugOut, err
}

func (s *shortURLs) BatchCreateShortURLs(
	ctx context.Context,
	userID string,
	items []models.BatchCreateShortURLs,
) ([]models.BatchCreateShortURLsResult, error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("%s: %w", storage.ErrMsgSomethingWentWrong, err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	itemsNum := len(items)
	if itemsNum == 0 {
		return nil, storage.ErrNoItemsInBatch
	}

	userIDs := make([]string, itemsNum)
	slugs := make([]string, itemsNum)
	originalURLs := make([]string, itemsNum)

	for idx, item := range items {
		userIDs[idx] = userID
		slugs[idx] = item.Slug
		originalURLs[idx] = item.OriginalURL
	}

	slugsOut, err := qtx.BatchCreateShortURLs(ctx,
		shorturls.BatchCreateShortURLsParams{
			UserIds:      userIDs,
			Slugs:        slugs,
			OriginalUrls: originalURLs,
		})
	if err != nil {
		err = fmt.Errorf("%s: %w", storage.ErrMsgSomethingWentWrong, err)
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		err = fmt.Errorf("%s: %w", storage.ErrMsgSomethingWentWrong, err)
		return nil, err
	}

	result := make([]models.BatchCreateShortURLsResult, itemsNum)
	for idx, slug := range slugsOut {
		result[idx].Slug = slug
		result[idx].ShortURL = s.baseURL + slug
	}

	return result, nil
}
