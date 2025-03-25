package services

import (
	"context"
	"fmt"
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

func (s *resolver) validateID(id string) error {
	const op = "service.resolver.validate_id"

	isValid := s.idRegexp.MatchString(id)
	if !isValid {
		return fmt.Errorf("%s - %w", op, ErrInvalidID)
	}

	return nil
}

func (s *resolver) Serve(ctx context.Context, id string) (string, error) {
	const op = "service.resolver.serve"

	err := s.validateID(id)
	if err != nil {
		return "", err
	}

	url, err := s.cache.GetURL(ctx, id)
	if err != nil {
		return "", err
	}

	if url == "" {
		return "", fmt.Errorf("%s - %w", op, ErrURLNotFound)
	}

	return url, err
}
