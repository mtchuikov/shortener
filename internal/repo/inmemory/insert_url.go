package inmemory

import (
	"context"
	"fmt"

	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/repo"
)

func (r *shortenURLs) InsertShortenURL(
	ctx context.Context,
	shortenID models.ShortenID,
	originalURL models.OriginalURL,
) (
	models.ShortenID,
	error,
) {
	r.rmu.Lock()
	defer r.rmu.Unlock()

	newRawShortenID := shortenID.String()
	rawOriginalURL := originalURL.String()

	rawShortenID, exists := r.shortenIDs[rawOriginalURL]
	if exists {
		shortenID, _ = models.NewShortenID(rawShortenID)
		return shortenID, repo.ErrOirignalURLAlreadyShortened
	}

	r.originalURLs[newRawShortenID] = rawOriginalURL
	r.shortenIDs[rawOriginalURL] = newRawShortenID

	err := r.backup(rawOriginalURL, newRawShortenID)
	if err != nil {
		return "", fmt.Errorf("%w: %s", repo.ErrUnexpectedError, err)
	}

	return shortenID, nil
}
