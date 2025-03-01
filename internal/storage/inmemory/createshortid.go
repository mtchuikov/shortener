package storage

import (
	"context"
)

func (s *Storage) CreateShortID(ctx context.Context, shortID, originalURL string) error {
	s.mu.Lock()
	s.originalURLs[shortID] = originalURL
	s.shortIDs[originalURL] = shortID
	s.mu.Unlock()

	return nil
}
