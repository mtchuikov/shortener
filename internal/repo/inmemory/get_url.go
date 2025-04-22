package inmemory

import (
	"context"

	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/repo"
)

func (r *shortenURLs) GetOriginalURL(
	ctx context.Context,
	shortenID models.ShortenID,
) (
	models.OriginalURL,
	error,
) {
	r.rmu.RLock()
	defer r.rmu.RUnlock()

	rawOriginalURL, exists := r.originalURLs[shortenID.String()]
	if !exists {
		return "", repo.ErrOriginalURLNotFound
	}

	originalURL, _ := models.NewOriginalURL(rawOriginalURL)
	return originalURL, nil
}
