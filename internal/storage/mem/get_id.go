package memstorage

import "context"

func (s *Storage) GetID(_ context.Context, url string) (string, error) {
	val := s.ids[url]
	return val, nil
}
