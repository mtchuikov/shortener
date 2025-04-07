package repository

import "context"

type Repo interface {
	CreateShortURL(ctx context.Context, originalURL, shortID string) error
	GetOriginalURL(ctx context.Context, shortID string) (string, error)
	GetShortID(ctx context.Context, originalURL string) (string, error)
}
