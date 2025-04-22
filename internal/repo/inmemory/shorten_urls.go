package inmemory

import (
	"os"
	"sync"
)

type shortenURLs struct {
	rmu          sync.RWMutex
	originalURLs map[string]string // id => url
	shortenIDs   map[string]string // url => id

	file *os.File
}

func New(file string) (*shortenURLs, error) {
	flags := os.O_CREATE | os.O_RDWR | os.O_APPEND
	mode := os.FileMode(0755)

	f, err := os.OpenFile(file, flags, mode)
	if err != nil {
		return nil, err
	}

	shortenURLs := &shortenURLs{
		rmu:          sync.RWMutex{},
		originalURLs: map[string]string{},
		shortenIDs:   map[string]string{},
		file:         f,
	}

	err = shortenURLs.restoreBackup()
	if err != nil {
		return nil, err
	}

	return shortenURLs, nil
}

func (r *shortenURLs) Close() {
	r.file.Close()
}
