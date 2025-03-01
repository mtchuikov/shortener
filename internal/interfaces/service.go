package interfaces

import "context"

type Service interface {
	CreateShortID(ctx context.Context, originalURL string) (string, int, error)
	ResolveShortID(ctx context.Context, shortID string) (string, int, error)
}
