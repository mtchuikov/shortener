package services

import (
	"context"
	"errors"
	"regexp"

	"github.com/mtchuikov/shortener/pkg/strgen"
)

type shortenerRepo interface {
	CreateShortURL(ctx context.Context, originalURL, shortID string) error
	GetShortID(ctx context.Context, originalURL string) (string, error)
}

type shortener struct {
	baseURL     string
	originalURL *regexp.Regexp
	shortIDGen  *strgen.Generator
	repo        shortenerRepo
}

func NewShortener(baseURL string, repo shortenerRepo) *shortener {
	return &shortener{
		baseURL:     baseURL,
		originalURL: regexp.MustCompile(`^(http://|https://)[a-zA-Z0-9]+([-.][a-zA-Z0-9]+)*\.[a-zA-Z]{2,}(:[0-9]{1,5})?(/.*)?$`),
		shortIDGen:  strgen.New(),
		repo:        repo,
	}
}

var ErrInvalidURL = errors.New("invalid url")

func (s *shortener) validateURL(url string) error {
	isValid := s.originalURL.MatchString(url)
	if !isValid {
		return ErrInvalidURL
	}

	return nil
}

var ErrFailedToCreateShortURL = errors.New("failed to create short url")

func (s *shortener) Serve(ctx context.Context, originalURL string) (string, error) {
	err := s.validateURL(originalURL)
	if err != nil {
		return "", err
	}

	shortID, err := s.repo.GetShortID(ctx, originalURL)
	if shortID == "" && err == nil {
		shortID = s.shortIDGen.Generate(8)
		err = s.repo.CreateShortURL(ctx, originalURL, shortID)
	}

	if err != nil {
		return "", ErrFailedToCreateShortURL
	}

	shortURL := s.baseURL + shortID
	return shortURL, nil
}
