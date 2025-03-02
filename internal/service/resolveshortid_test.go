package service

import (
	"context"
	"net/http"
	"testing"

	storage "github.com/mtchuikov/shortener/internal/storage/inmemory"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestIsValidShortID(t *testing.T) {
	svc := Service{}

	validShortIDs := []string{
		"AbCd1234", "ABCDEFGH",
		"aBcDeF01", "12345678",
	}

	for _, shortID := range validShortIDs {
		require.True(t, svc.isValidShortID(shortID), "expected valid short id: %q", shortID)
	}

	invalidShortIDs := []string{
		"avc", "$BCDEFGH",
		"aBcDeF!1", "",
	}

	for _, shortID := range invalidShortIDs {
		require.False(t, svc.isValidShortID(shortID), "expected invalid short id: %q", shortID)
	}
}

func TestResolveShortID_Success(t *testing.T) {
	storage := storage.New()
	service := New(log.Logger, storage)

	ctx := context.Background()
	mockOriginalURL := "http://example.com"
	shortID := "aBcDeF01"

	err := storage.CreateShortID(ctx, shortID, mockOriginalURL)
	require.NoError(t, err, "failed to create short id in storage")

	originalURL, code, err := service.ResolveShortID(ctx, shortID)
	require.Equal(t, mockOriginalURL, originalURL, "expected original url")
	require.Equal(t, 0, code, "expected 0 status code for success")
	require.Nil(t, err, "expected no error")
}

func TestResolveShortID_InvalidShortID(t *testing.T) {
	service := New(log.Logger, nil)

	ctx := context.Background()
	invalidShortID := "#dQ41f9!"

	originalURL, code, err := service.ResolveShortID(ctx, invalidShortID)
	require.Empty(t, originalURL, "expected empty original url")
	require.Equal(t, http.StatusBadRequest, code, "expected 400 status code for invalid short id")
	require.ErrorIs(t, err, ErrInvalidShortID, "expected invalid short id error")
}

func TestResolveShortID_NoSuchShortID(t *testing.T) {
	storage := storage.New()
	service := New(log.Logger, storage)

	ctx := context.Background()
	shortID := "aBcDeF01"

	originalURL, code, err := service.ResolveShortID(ctx, shortID)
	require.Empty(t, originalURL, "expected empty original url")
	require.Equal(t, http.StatusNotFound, code, "expected 404 status code for non-existent short id")
	require.ErrorIs(t, err, ErrNoSuchShortID, "expected no such short id error")
}
