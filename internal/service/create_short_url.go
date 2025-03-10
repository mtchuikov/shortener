package service

import (
	"context"
	"net/url"

	"github.com/mtchuikov/shortener/pkg/randtools"
)

func isValidOriginalURL(originalURL string) error {
	url, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return ErrInvalidOriginalURL
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return ErrInvalidOriginalURL
	}

	if url.Host == "" {
		return ErrInvalidOriginalURL
	}

	return nil
}

const shortURLLen = 8

func (s *Service) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	err := isValidOriginalURL(originalURL)
	if err != nil {
		return "", err
	}

	shortURL, err := s.repo.GetShortURL(ctx, originalURL)
	if err != nil {
		return "", ErrInternalError
	}

	if shortURL == "" {
		shortURL = randtools.DefaultGenerateString(shortURLLen)
		err := s.repo.CreateShortURL(ctx, originalURL, shortURL)
		if err != nil {
			return "", ErrInternalError
		}
	}

	return s.baseURL + shortURL, nil
}
