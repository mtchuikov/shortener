package services

import (
	"context"
	"errors"
	"regexp"

	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/repo"
	"github.com/mtchuikov/shortener/pkg/strgen"
)

type shortenerRepo interface {
	CreateShortURL(ctx context.Context, originalURL, shortID string) error
	BatchCreateShortURLs(ctx context.Context, urlsToShort models.URLsToShort) error
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

// TODO: rewrite this function cause it's implemented bad :(
func (s *shortenerService) ServeBatch(
	ctx context.Context,
	urlsToShort models.URLsToShort,
	numUrlsToShort int,
) (
	models.ShortenURLs,
	error,
) {
	for i := range numUrlsToShort {
		urlsToShort[i].ShortID = s.shortIDGen.Generate(s.shortIDLen)
	}

	err := s.repo.BatchCreateShortURLs(ctx, urlsToShort)
	if err != nil {
		return nil, err
	}

	shortenURLs := make(models.ShortenURLs, numUrlsToShort)
	for i := range numUrlsToShort {
		shortenURLs[i].CorrelationID = urlsToShort[i].CorrelationID
		shortenURLs[i].ShortURL = s.baseURL + urlsToShort[i].ShortID
	}

	return shortenURLs, nil
}
