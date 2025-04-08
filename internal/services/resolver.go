package services

import (
	"context"
	"regexp"
)

type resolverRepo interface {
	GetOriginalURL(ctx context.Context, shortID string) (string, error)
}

type resolverService struct {
	repo resolverRepo
}

func NewResolver(repo resolverRepo) *resolverService {
	return &resolverService{repo}
}

var shortIDRegexp = regexp.MustCompile(`^[A-Za-z0-9]{8}$`)

func (s *resolverService) validateShortID(id string) error {
	valid := shortIDRegexp.MatchString(id)
	if !valid {
		return ErrInvalidShortID
	}

	return nil
}

func (s *resolverService) Serve(ctx context.Context, shortID string) (string, error) {
	err := s.validateShortID(shortID)
	if err != nil {
		return "", err
	}

	return s.repo.GetOriginalURL(ctx, shortID)
}
