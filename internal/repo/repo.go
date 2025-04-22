package repo

import (
	"context"

	"github.com/mtchuikov/shortener/internal/models"
)

type IShortenURLs interface {
	InsertShortenURL(context.Context, models.ShortenID, models.OriginalURL) (
		models.ShortenID, error,
	)
	GetOriginalURL(context.Context, models.ShortenID) (models.OriginalURL, error)
}
