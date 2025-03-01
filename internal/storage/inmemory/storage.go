package storage

import "sync"

type Storage struct {
	mu           sync.Mutex
	shortIDs     map[string]string // short id => original url
	originalURLs map[string]string // original url => short id
}

func New() *Storage {
	return &Storage{
		mu:           sync.Mutex{},
		shortIDs:     map[string]string{},
		originalURLs: map[string]string{},
	}
}
