package service

import (
	"context"
	"net/http"
	"net/url"

	"github.com/mtchuikov/shortener/pkg/randtools"
)

func (s *Service) isValidURL(originalURL string) bool {
	url, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return false
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return false
	}

	if url.Host == "" {
		return false
	}

	return true
}

const shortIDLen = 8

func (s *Service) CreateShortID(ctx context.Context, originalURL string) (string, int, error) {
	valid := s.isValidURL(originalURL)
	if !valid {
		return "", http.StatusBadRequest, ErrInvalidURL
	}

	shortID, err := s.storage.GetShortID(ctx, originalURL)
	if err != nil {
		errMsg := "service: failed to check if original url exists"
		s.logger.Err(err).Msg(errMsg)

		code := http.StatusInternalServerError
		return "", code, ErrSomethidWentWrong
	}

	if shortID == "" {
		shortID = randtools.DefaultGenerateString(shortIDLen)
		err := s.storage.CreateShortID(ctx, shortID, originalURL)
		if err != nil {
			errMsg := "service: failed to write short id to db"
			s.logger.Err(err).Msg(errMsg)

			code := http.StatusInternalServerError
			return "", code, ErrSomethidWentWrong
		}
	}

	return shortID, 0, nil
}
