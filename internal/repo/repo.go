package repo

import "context"

type IRepo interface {
	CreateShortURL(ctx context.Context, originalURL, shortURL string) error
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	GetShortURL(ctx context.Context, originalURL string) (string, error)
}
