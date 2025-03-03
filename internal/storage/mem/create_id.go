package memstorage

import "context"

func (s *Storage) CreateID(_ context.Context, id, url string) error {
	s.mu.Lock()
	s.ids[url] = id
	s.urls[id] = url
	s.mu.Unlock()

	return nil
}
