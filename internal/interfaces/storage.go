package interfaces

import "context"

type Storage interface {
	CreateID(ctx context.Context, id, url string) error
	GetURL(ctx context.Context, id string) (string, error)
	GetID(ctx context.Context, url string) (string, error)
}
