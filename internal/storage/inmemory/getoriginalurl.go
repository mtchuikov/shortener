package storage

import (
	"context"
)

func (s *Storage) GetOriginalURL(_ context.Context, shortID string) (string, error) {
	val, ok := s.originalURLs[shortID]
	if !ok {
		return "", nil
	}

	return val, nil
}
