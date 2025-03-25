package inmemory

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
)

type cache struct {
	rmu  sync.RWMutex
	urls map[string]string // id => url
	ids  map[string]string // url => id

	// backup
	file *os.File
}

func New(backupFile string) (*cache, error) {
	const op = "cache.inmemory.cache.new"

	flags := os.O_CREATE | os.O_RDWR | os.O_APPEND
	mode := os.FileMode(0755)

	f, err := os.OpenFile(backupFile, flags, mode)
	if err != nil {
		err = errors.Unwrap(err)
		return nil, fmt.Errorf(
			"%s - %w: %v",
			op, ErrFailedToOpenBackupFile, err,
		)
	}

	cache := &cache{
		rmu:  sync.RWMutex{},
		urls: map[string]string{},
		ids:  map[string]string{},
		file: f,
	}

	err = cache.restoreBackup()
	if err != nil {
		return nil, err
	}

	return cache, nil
}
func (c *cache) Close(_ context.Context) error {
	return nil
}

func (c *cache) CreateShortURL(_ context.Context, url, id string) error {
	c.rmu.Lock()
	defer c.rmu.Unlock()

	c.urls[id] = url
	c.ids[url] = id

	return c.backup(url, id)
}

func (c *cache) GetURL(_ context.Context, id string) (string, error) {
	c.rmu.RLock()
	defer c.rmu.RUnlock()
	return c.urls[id], nil
}

func (c *cache) GetID(_ context.Context, url string) (string, error) {
	c.rmu.RLock()
	defer c.rmu.RUnlock()
	return c.ids[url], nil
}
