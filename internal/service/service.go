package service

import (
	"github.com/mtchuikov/shortener/internal/interfaces"
	"github.com/rs/zerolog"
)

type Service struct {
	logger  zerolog.Logger
	storage interfaces.Repo
}

func New(logger zerolog.Logger, storage interfaces.Repo) *Service {
	return &Service{
		logger:  logger,
		storage: storage,
	}
}
