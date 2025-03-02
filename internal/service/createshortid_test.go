package service

import (
	"context"
	"net/http"
	"testing"

	storage "github.com/mtchuikov/shortener/internal/storage/inmemory"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestIsValidURL(t *testing.T) {
	service := Service{}

	validURLs := []string{
		"http://example.com",
		"http://yandex.ru",
		"http://sub.domain.ru",
		"https://example.com/path?query=1",
	}

	for _, u := range validURLs {
		require.True(t, service.isValidURL(u), "expected valid url: %q", u)
	}

	invalidURLs := []string{
		"ftp://example.com",
		"",
		"http://",
		"yandex.ru",
	}

	for _, u := range invalidURLs {
		require.False(t, service.isValidURL(u), "expected invalid url: %q", u)
	}
}

func TestCreateShortID_Success(t *testing.T) {
	storage := storage.New()
	service := New(log.Logger, storage)

	ctx := context.Background()
	originalURL := "http://example.com"

	shortID, code, err := service.CreateShortID(ctx, originalURL)
	require.True(t, service.isValidShortID(shortID), "expected valid short id: %q", shortID)
	require.Zero(t, code, "expected 0 status code for success")
	require.NoError(t, err, "expected no error")
}

func TestCreateShortID_InvalidURL(t *testing.T) {
	service := New(log.Logger, nil)

	ctx := context.Background()
	invalidURL := "not url"

	shortID, code, err := service.CreateShortID(ctx, invalidURL)
	require.Empty(t, shortID, "expected empty short id")
	require.Equal(t, http.StatusBadRequest, code, "expected bad request status code")
	require.ErrorIs(t, err, ErrInvalidURL, "expected invalid url error")
}
