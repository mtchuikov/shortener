package inmemory

import (
	"context"
	"sync"
)

type cache struct {
	mu   sync.Mutex
	urls map[string]string // id => url
	ids  map[string]string // url => id
}

func New() *cache {
	return &cache{
		mu:   sync.Mutex{},
		urls: map[string]string{},
		ids:  map[string]string{},
	}
}

func (c *cache) CreateShortURL(_ context.Context, url, id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.urls[id] = url
	c.ids[url] = id

	return nil
}

func (c *cache) GetURL(_ context.Context, id string) (string, error) {
	return c.urls[id], nil
}

func (c *cache) GetID(_ context.Context, url string) (string, error) {
	return c.ids[url], nil
}
