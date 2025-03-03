package interfaces

import "context"

type Service interface {
	ShortURL(ctx context.Context, url string) (string, error)
	ResolveShortID(ctx context.Context, id string) (string, error)
}
