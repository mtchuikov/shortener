package service

import (
	"strings"

	"github.com/mtchuikov/shortener/internal/interfaces"
	"github.com/rs/zerolog"
)

type Service struct {
	baseURL string
	logger  zerolog.Logger
	storage interfaces.Storage
}

func New(logger zerolog.Logger, baseURL string, storage interfaces.Storage) *Service {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}

	return &Service{
		storage: storage,
		logger:  logger,
		baseURL: baseURL,
	}
}
