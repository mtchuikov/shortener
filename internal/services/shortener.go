package services

import (
	"context"
	"errors"
	"regexp"

	"github.com/mtchuikov/shortener/internal/repo"
	"github.com/mtchuikov/shortener/pkg/strgen"
)

type shortenerRepo interface {
	CreateShortURL(ctx context.Context, originalURL, shortID string) error
	GetShortID(ctx context.Context, originalURL string) (string, error)
}

type shortenerService struct {
	baseURL    string
	shortIDLen int
	shortIDGen *strgen.Generator
	repo       shortenerRepo
}

func NewShortener(baseURL string, repo shortenerRepo) *shortenerService {
	return &shortenerService{
		baseURL:    baseURL,
		shortIDLen: 8,
		shortIDGen: strgen.New(),
		repo:       repo,
	}
}

var originalURLRegexp = regexp.MustCompile(`^(http://|https://)[a-zA-Z0-9]+([-.][a-zA-Z0-9]+)*\.[a-zA-Z]{2,}(:[0-9]{1,5})?(/.*)?$`)

func (s *shortenerService) validateOriginalURL(u string) error {
	valid := originalURLRegexp.MatchString(u)
	if !valid {
		return ErrInvalidOriginalURL
	}

	return nil
}

func (s *shortenerService) Serve(ctx context.Context, originalURL string) (string, error) {
	err := s.validateOriginalURL(originalURL)
	if err != nil {
		return "", err
	}

	shortID, err := s.repo.GetShortID(ctx, originalURL)
	if err == nil {
		shortURL := s.baseURL + shortID
		return shortURL, ErrOriginalURLAlreadyShorten
	}

	if !errors.Is(err, repo.ErrShortIDNotFound) {
		return "", err
	}

	shortID = s.shortIDGen.Generate(s.shortIDLen)
	err = s.repo.CreateShortURL(ctx, originalURL, shortID)
	if err != nil {
		return "", err
	}

	shortURL := s.baseURL + shortID
	return shortURL, nil
}
