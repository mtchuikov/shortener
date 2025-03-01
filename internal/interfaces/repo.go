package interfaces

import "context"

type Repo interface {
	CreateShortID(ctx context.Context, shortID, originalURL string) error
	GetOriginalURL(ctx context.Context, shortID string) (string, error)
	GetShortID(ctx context.Context, originalURL string) (string, error)
}
