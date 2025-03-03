package service

import (
	"context"
	"regexp"
	"testing"

	memstorage "github.com/mtchuikov/shortener/internal/storage/mem"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestShortURL_Success(t *testing.T) {
	rootCtx := context.Background()
	storage := memstorage.New()

	baseURL := "http://localhost:3214/api/"
	service := New(log.Logger, baseURL, storage)

	urls := []string{
		"https://example.com", "https://example1.com",
		"https://example.com/page", "https://example.com?limit=1",
	}

	for _, url := range urls {
		shortURL, err := service.ShortURL(rootCtx, url)
		require.NoError(t, err, "expected no error when calling service")

		var validShortURL = regexp.MustCompile(`^http://localhost:3214/api/[A-Za-z0-9]{8}$`)
		valid := validShortURL.MatchString(shortURL)
		require.True(t, valid, "expected valid short url")
	}
}

func TestCreateShortID_InvalidURL(t *testing.T) {
	rootCtx := context.Background()
	storage := memstorage.New()
	service := New(log.Logger, "", storage)

	urls := []string{
		"ftp://example.com", "",
		"http://", "example.com",
	}

	for _, url := range urls {
		shortURL, err := service.ShortURL(rootCtx, url)
		require.ErrorIs(t, err, ErrInvalidURL, "expected invalid url error")
		require.Empty(t, shortURL, "expected emptry short url")
	}

}
