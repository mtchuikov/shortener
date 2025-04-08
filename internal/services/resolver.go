package services

import (
	"context"
)

type resolverRepo interface {
	GetOriginalURL(ctx context.Context, shortID string) (string, error)
}

type resolverService struct {
	repo resolverRepo
}

func NewResolver(repo resolverRepo) *resolverService {
	return &resolverService{repo}
}

func (s *resolverService) Serve(ctx context.Context, shortID string) (string, error) {
	return s.repo.GetOriginalURL(ctx, shortID)
}
