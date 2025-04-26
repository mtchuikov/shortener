package inmemory

import (
	"context"
	"fmt"

	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/storage"
)

func (s *shortURLs) CreateShortURL(
	ctx context.Context,
	userID,
	slug,
	originalURL string,
) (string, error) {
	s.rmu.Lock()
	defer s.rmu.Unlock()

	slugOut, exists := s.slugs[originalURL]
	if exists {
		return s.baseURL + slugOut, storage.ErrOriginalURLAlreadyExists
	}

	s.slugs[originalURL] = slug
	s.originalURLs[slug] = originalURL

	err := s.backup(slug, originalURL)
	if err != nil {
		delete(s.slugs, originalURL)
		delete(s.originalURLs, slug)

		return "", err
	}

	return s.baseURL + slug, nil
}

func (s *shortURLs) BatchCreateShortURLs(
	ctx context.Context,
	userID string,
	items []models.BatchCreateShortURLs,
) ([]models.BatchCreateShortURLsResult, error) {
	err := fmt.Errorf("%s: unimplemented", storage.ErrMsgSomethingWentWrong)
	return nil, err
}
