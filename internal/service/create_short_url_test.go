package service

import (
	"context"
	"regexp"
	"testing"

	"github.com/mtchuikov/shortener/internal/repo/inmemory"
	"github.com/stretchr/testify/require"
)

func TestCreateShortURL_Success(t *testing.T) {
	rootCtx := context.Background()
	baseURL := "http://localhost:3214/api/"
	service := New(baseURL, inmemory.New())

	urls := []string{
		"https://example.com", "https://example1.com",
		"https://example.com/page", "https://example.com?limit=1",
	}

	for _, url := range urls {
		shortURL, err := service.CreateShortURL(rootCtx, url)
		require.NoError(t, err, "expected no error when calling service")

		var validShortURL = regexp.MustCompile(`^http://localhost:3214/api/[A-Za-z0-9]{8}$`)
		valid := validShortURL.MatchString(shortURL)
		require.True(t, valid, "expected valid short url")
	}
}

func TestCreateShortURL_InvalidShortID(t *testing.T) {
	rootCtx := context.Background()
	baseURL := "http://localhost:3214/api/"
	service := New(baseURL, inmemory.New())

	urls := []string{
		"ftp://example.com", "",
		"http://", "example.com",
	}

	for _, url := range urls {
		shortURL, err := service.CreateShortURL(rootCtx, url)
		require.ErrorIs(t, err, ErrInvalidOriginalURL, "expected invalid url error")
		require.Empty(t, shortURL, "expected emptry short url")
	}

}
