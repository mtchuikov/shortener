package inmemory

import (
	"context"
	"fmt"
	"os"
	"sync"
)

type repo struct {
	rmu          sync.RWMutex
	originalURLs map[string]string // id => url
	shortIDs     map[string]string // url => id

	file *os.File
}

const wFailedToSetupInMemoryRepo = "failed to setup inmemory repo: %w"

func New(backupFile string) (*repo, error) {
	flags := os.O_CREATE | os.O_RDWR | os.O_APPEND
	mode := os.FileMode(0755)

	f, err := os.OpenFile(backupFile, flags, mode)
	if err != nil {
		return nil, fmt.Errorf(wFailedToSetupInMemoryRepo, err)
	}

	inmemory := &repo{
		rmu:          sync.RWMutex{},
		originalURLs: map[string]string{},
		shortIDs:     map[string]string{},
		file:         f,
	}

	err = inmemory.restoreBackup()
	if err != nil {
		return nil, fmt.Errorf(wFailedToSetupInMemoryRepo, err)
	}

	return inmemory, nil
}

func (r *repo) CreateShortURL(ctx context.Context, url, id string) error {
	r.rmu.Lock()
	defer r.rmu.Unlock()

	r.originalURLs[id] = url
	r.shortIDs[url] = id

	return r.backup(url, id)
}

func (r *repo) GetOriginalURL(ctx context.Context, shortID string) (string, error) {
	r.rmu.RLock()
	defer r.rmu.RUnlock()
	return r.originalURLs[shortID], nil
}

func (r *repo) GetShortID(ctx context.Context, originalURL string) (string, error) {
	r.rmu.RLock()
	defer r.rmu.RUnlock()
	return r.shortIDs[originalURL], nil
}
