package memstorage

import "sync"

type Storage struct {
	mu   sync.Mutex
	ids  map[string]string // url => id
	urls map[string]string // id => url
}

func New() *Storage {
	return &Storage{
		mu:   sync.Mutex{},
		ids:  map[string]string{},
		urls: map[string]string{},
	}
}
