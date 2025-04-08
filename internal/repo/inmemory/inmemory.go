package inmemory

import (
	"context"
	"os"
	"sync"

	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/repo"
)

type inmemory struct {
	rmu          sync.RWMutex
	originalURLs map[string]string // id => url
	shortIDs     map[string]string // url => id

	file *os.File
}

func New(backupFile string) (*inmemory, error) {
	flags := os.O_CREATE | os.O_RDWR | os.O_APPEND
	mode := os.FileMode(0755)

	f, err := os.OpenFile(backupFile, flags, mode)
	if err != nil {
		return nil, err
	}

	inmemory := &inmemory{
		rmu:          sync.RWMutex{},
		originalURLs: map[string]string{},
		shortIDs:     map[string]string{},
		file:         f,
	}

	err = inmemory.restoreBackup()
	if err != nil {
		return nil, err
	}

	return inmemory, nil
}

func (r *inmemory) CreateShortURL(ctx context.Context, originalURL, shortID string) error {
	r.rmu.Lock()
	defer r.rmu.Unlock()

	r.originalURLs[shortID] = originalURL
	r.shortIDs[originalURL] = shortID

	return r.backup(originalURL, shortID)
}

func (r *inmemory) BatchCreateShortURLs(ctx context.Context, urlsToShort models.URLsToShort) error {
	r.rmu.Lock()
	defer r.rmu.Unlock()

	for _, urlToShot := range urlsToShort {
		r.originalURLs[urlToShot.CorrelationID] = urlToShot.OriginalURL
		r.shortIDs[urlToShot.OriginalURL] = urlToShot.CorrelationID

		err := r.backup(urlToShot.OriginalURL, urlToShot.CorrelationID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *inmemory) GetOriginalURL(ctx context.Context, shortID string) (string, error) {
	r.rmu.RLock()
	defer r.rmu.RUnlock()

	originalURL, exists := r.originalURLs[shortID]
	if !exists {
		return "", repo.ErrOriginalURLNotFound
	}

	return originalURL, nil
}

func (r *inmemory) GetShortID(ctx context.Context, originalURL string) (string, error) {
	r.rmu.RLock()
	defer r.rmu.RUnlock()

	shortID, exists := r.shortIDs[originalURL]
	if !exists {
		return "", repo.ErrShortIDNotFound
	}

	return shortID, nil
}

func (r *inmemory) Close() {
	r.file.Close()
}
