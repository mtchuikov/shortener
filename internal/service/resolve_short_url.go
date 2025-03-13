package service

import (
	"context"
	"regexp"
)

var shortURLRegex = regexp.MustCompile(`^[A-Za-z0-9]{8}$`)

func isValidShortURL(shortURL string) error {
	match := shortURLRegex.MatchString(shortURL)
	if !match {
		return ErrInvalidShortURL
	}

	return nil
}

func (s *Service) ResolveShortURL(ctx context.Context, shortURL string) (string, error) {
	err := isValidShortURL(shortURL)
	if err != nil {
		return "", err
	}

	originalURL, err := s.repo.GetOriginalURL(ctx, shortURL)
	if err != nil {
		return "", ErrInternalError
	}

	if originalURL == "" {
		return "", ErrShortURLNotFound
	}

	return originalURL, nil
}
