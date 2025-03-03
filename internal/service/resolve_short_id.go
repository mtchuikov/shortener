package service

import (
	"context"
	"regexp"
)

var idRegex = regexp.MustCompile(`^[A-Za-z0-9]{8}$`)

func validateID(id string) error {
	match := idRegex.MatchString(id)
	if !match {
		return ErrInvalidID
	}

	return nil
}

func (s *Service) ResolveShortID(ctx context.Context, id string) (string, error) {
	err := validateID(id)
	if err != nil {
		return "", err
	}

	url, err := s.storage.GetURL(ctx, id)
	if err != nil {
		return "", ErrInternalError
	}

	if url == "" {
		return "", ErrIDNotFound
	}

	return url, nil
}
