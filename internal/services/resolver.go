package services

import (
	"context"
	"errors"
	"regexp"
)

type resolverRepo interface {
	GetOriginalURL(ctx context.Context, shortID string) (string, error)
}

type resolver struct {
	shortIDRegexp *regexp.Regexp
	repo          resolverRepo
}

func NewResolver(repo resolverRepo) *resolver {
	return &resolver{
		shortIDRegexp: regexp.MustCompile(`^[A-Za-z0-9]{8}$`),
		repo:          repo,
	}
}

var ErrInvalidShortID = errors.New("invalid short id")

func (s *resolver) validateID(id string) error {
	isValid := s.shortIDRegexp.MatchString(id)
	if !isValid {
		return ErrInvalidShortID
	}

	return nil
}

var (
	ErrFailedToGetOriginalURL = errors.New("failed to get original url")
	ErrOriginalURLNotFound    = errors.New("original url not found")
)

func (s *resolver) Serve(ctx context.Context, shortID string) (string, error) {
	err := s.validateID(shortID)
	if err != nil {
		return "", err
	}

	originalURL, err := s.repo.GetOriginalURL(ctx, shortID)
	if err != nil {
		return "", ErrFailedToGetOriginalURL
	}

	if originalURL == "" {
		return "", ErrOriginalURLNotFound
	}

	return originalURL, nil
}
