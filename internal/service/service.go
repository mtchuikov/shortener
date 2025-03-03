package service

import (
	"github.com/mtchuikov/shortener/internal/interfaces"
	"github.com/rs/zerolog"
)

type Service struct {
	baseURL string
	logger  zerolog.Logger
	storage interfaces.Storage
}

func New(logger zerolog.Logger, baseURL string, storage interfaces.Storage) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
		baseURL: baseURL,
	}
}
