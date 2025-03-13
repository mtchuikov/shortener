package inmemory

import "sync"

type Repo struct {
	mu           sync.Mutex
	originalURLs map[string]string // short url => original url
	shortURLs    map[string]string // original url => short url
}

func New() *Repo {
	return &Repo{
		mu:           sync.Mutex{},
		originalURLs: map[string]string{},
		shortURLs:    map[string]string{},
	}
}
