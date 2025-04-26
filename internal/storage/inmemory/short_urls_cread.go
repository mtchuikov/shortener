package inmemory

import (
	"context"
	"fmt"

	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/storage"
)

func (s *shortURLs) GetOriginalURLBySlug(ctx context.Context, slug string) (string, bool, error) {
	s.rmu.RLock()
	defer s.rmu.RUnlock()

	originalURL, exists := s.originalURLs[slug]
	if !exists {
		return "", false, storage.ErrOriginalURLNotFound
	}

	return originalURL, false, nil
}

func (s *shortURLs) ListShortURLsByUser(
	ctx context.Context,
	userID string,
) ([]models.ListShortURLsByUserResult, error) {
	err := fmt.Errorf("%s: unimplemented", storage.ErrMsgSomethingWentWrong)
	return nil, err
}
