package service

import (
	"context"
	"net/http"
	"regexp"
)

var shortIDRegex = regexp.MustCompile(`^[A-Za-z0-9]{8}$`)

func (s *Service) isValidShortID(shortID string) bool {
	return shortIDRegex.MatchString(shortID)
}

func (s *Service) ResolveShortID(ctx context.Context, shortID string) (string, int, error) {
	valid := s.isValidShortID(shortID)
	if !valid {
		return "", http.StatusBadRequest, ErrInvalidShortID
	}

	originalURL, err := s.storage.GetOriginalURL(ctx, shortID)
	if err != nil {
		errMsg := "service: failed to get original url"
		s.logger.Err(err).Msg(errMsg)

		code := http.StatusInternalServerError
		return "", code, ErrSomethidWentWrong
	}

	if originalURL == "" {
		return "", http.StatusNotFound, ErrNoSuchShortID
	}

	return originalURL, 0, nil
}
