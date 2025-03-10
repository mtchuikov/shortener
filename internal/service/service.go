package service

import (
	"context"

	"github.com/mtchuikov/shortener/internal/config"
	"github.com/mtchuikov/shortener/internal/repo"
)

type IService interface {
	CreateShortURL(ctx context.Context, originalURL string) (string, error)
	ResolveShortURL(ctx context.Context, shortURL string) (string, error)
}

type Service struct {
	repo    repo.IRepo
	baseURL string
}

func New(config config.Config, repo repo.IRepo) *Service {
	return &Service{
		repo:    repo,
		baseURL: config.BaseURL,
	}
}
