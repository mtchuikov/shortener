package service

import (
	"context"
	"net/url"

	"github.com/mtchuikov/shortener/pkg/randtools"
)

func validateURL(u string) error {
	url, err := url.ParseRequestURI(u)
	if err != nil {
		return ErrInvalidURL
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return ErrInvalidURL
	}

	if url.Host == "" {
		return ErrInvalidURL
	}

	return nil
}

var idLen = 8

func (s *Service) ShortURL(ctx context.Context, url string) (string, error) {
	err := validateURL(url)
	if err != nil {
		return "", err
	}

	id, err := s.storage.GetID(ctx, url)
	if err != nil {
		return "", ErrInternalError
	}

	if id == "" {
		id = randtools.DefaultGenerateString(idLen)
		err = s.storage.CreateID(ctx, id, url)
		if err != nil {
			return "", ErrInternalError
		}
	}

	return s.baseURL + id, nil
}
