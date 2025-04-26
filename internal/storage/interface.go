package storage

import (
	"context"

	"github.com/mtchuikov/shortener/internal/models"
)

type ShortURLsStorage interface {
	CreateShortURL(ctx context.Context, userID, slug, originalURL string) (string, error)
	BatchCreateShortURLs(ctx context.Context, userID string, items []models.BatchCreateShortURLs) ([]models.BatchCreateShortURLsResult, error)
	GetOriginalURLBySlug(ctx context.Context, slug string) (string, bool, error)
	ListShortURLsByUser(ctx context.Context, userID string) ([]models.ListShortURLsByUserResult, error)
	MarkShortURLsAsActive(ctx context.Context, userID string, slugs []string) error
	MarkShortURLsAsDeleted(ctx context.Context, userID string, slugs []string) error
}
