package memstorage

import "context"

func (s *Storage) GetURL(_ context.Context, id string) (string, error) {
	val := s.urls[id]
	return val, nil
}
