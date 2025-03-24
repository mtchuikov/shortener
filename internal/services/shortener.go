package services

import (
	"context"
	"errors"
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
		urlRegexp: regexp.MustCompile(`^(http|https)://[^:/\s]+`),
		idgen:     randtools.NewStringGenerator(),
		cache:     cache,
	}
}

var ErrInvalidURL = errors.New("invalid url")

func (s *shortener) validateURL(url string) error {
	isValid := s.urlRegexp.MatchString(url)
	if !isValid {
		return ErrInvalidURL
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
