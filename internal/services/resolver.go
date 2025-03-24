package services

import (
	"context"
	"errors"
	"regexp"
)

type resolverCache interface {
	GetURL(ctx context.Context, id string) (string, error)
}

type resolver struct {
	baseURL  string
	idRegexp *regexp.Regexp
	cache    resolverCache
}

func NewResolver(baseURL string, cache resolverCache) *resolver {
	return &resolver{
		baseURL:  baseURL,
		idRegexp: regexp.MustCompile(`^[A-Za-z0-9]{8}$`),
		cache:    cache,
	}
}

var ErrInvalidID = errors.New("invalid id")

func (s *resolver) validateID(id string) error {
	isValid := s.idRegexp.MatchString(id)
	if !isValid {
		return ErrInvalidID
	}

	return nil
}

var ErrIDNotFound = errors.New("id not found")

func (s *resolver) Serve(ctx context.Context, id string) (string, error) {
	err := s.validateID(id)
	if err != nil {
		return "", err
	}

	url, err := s.cache.GetURL(ctx, id)
	if err != nil {
		return "", err
	}

	if url == "" {
		return "", ErrIDNotFound
	}

	return url, err
}
