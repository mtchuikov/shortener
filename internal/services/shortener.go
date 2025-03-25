package services

import (
	"context"
	"fmt"
	"regexp"

	"github.com/mtchuikov/shortener/pkg/randtools"
)

type shortenerCache interface {
	CreateShortURL(ctx context.Context, url, id string) error
	GetID(ctx context.Context, url string) (string, error)
}

type shortener struct {
	baseURL   string
	urlRegexp *regexp.Regexp
	idgen     *randtools.StringGenerator
	cache     shortenerCache
}

func NewShortener(baseURL string, cache shortenerCache) *shortener {
	return &shortener{
		baseURL:   baseURL,
		urlRegexp: regexp.MustCompile(`^(http://|https://)[a-zA-Z0-9]+([-.][a-zA-Z0-9]+)*\.[a-zA-Z]{2,}(:[0-9]{1,5})?(/.*)?$`),
		idgen:     randtools.NewStringGenerator(),
		cache:     cache,
	}
}

func (s *shortener) validateURL(url string) error {
	const op = "service.shortener.validate_url"

	isValid := s.urlRegexp.MatchString(url)
	if !isValid {
		return fmt.Errorf("%s - %w", op, ErrInvalidURL)
	}

	return nil
}

func (s *shortener) Serve(ctx context.Context, url string) (string, error) {
	err := s.validateURL(url)
	if err != nil {
		return "", err
	}

	id, err := s.cache.GetID(ctx, url)
	if id == "" && err == nil {
		const idLen = 8
		id = s.idgen.Generate(idLen)

		err = s.cache.CreateShortURL(ctx, url, id)
	}

	if err != nil {
		return "", err
	}

	shortURL := s.baseURL + id
	return shortURL, nil
}
