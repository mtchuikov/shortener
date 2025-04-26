package inmemory

import (
	"os"
	"sync"

	"github.com/mtchuikov/shortener/internal/storage"
)

var _ storage.ShortURLsStorage = (*shortURLs)(nil)

type shortURLs struct {
	rmu          sync.RWMutex
	slugs        map[string]string
	originalURLs map[string]string
	backupFile   *os.File
	baseURL      string
}

func NewShortURLs(file string, baseURL string) (*shortURLs, error) {
	flags := os.O_CREATE | os.O_RDWR | os.O_APPEND
	mode := os.FileMode(0755)

	backupFile, err := os.OpenFile(file, flags, mode)
	if err != nil {
		return nil, err
	}

	shortURLs := &shortURLs{
		rmu:          sync.RWMutex{},
		slugs:        map[string]string{},
		originalURLs: map[string]string{},
		backupFile:   backupFile,
		baseURL:      baseURL,
	}

	err = shortURLs.restoreBackup()
	if err != nil {
		return nil, err
	}

	return shortURLs, nil
}

func (s *shortURLs) Close() {
	s.backupFile.Close()
}
