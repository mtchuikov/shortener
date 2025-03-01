package storage

import "context"

func (s *Storage) GetShortID(_ context.Context, originalURL string) (string, error) {
	val, ok := s.shortIDs[originalURL]
	if !ok {
		return "", nil
	}

	return val, nil
}
