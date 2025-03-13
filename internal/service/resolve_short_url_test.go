package service

import (
	"context"
	"testing"

	"github.com/mtchuikov/shortener/internal/repo/inmemory"
	"github.com/stretchr/testify/require"
)

func TestResolveShortURL_Success(t *testing.T) {
	rootCtx := context.Background()
	repo := inmemory.New()
	service := New("", repo)

	mockOriginalURL := "http://example.com"
	mockShortURL := "Stb1end7"

	err := repo.CreateShortURL(rootCtx, mockOriginalURL, mockShortURL)
	require.NoError(t, err, "expected no error when adding short url to repo")

	originalURL, err := service.ResolveShortURL(rootCtx, mockShortURL)
	require.NoError(t, err, "expected no error when resolving short url")

	require.Equal(t, mockOriginalURL, originalURL, "expected to get same original url")
}

func TestResolveShortURL_InvalidShortURL(t *testing.T) {
	rootCtx := context.Background()
	repo := inmemory.New()
	service := New("", repo)

	mockOriginalURL := "http://example.com"
	mockShortURL := ""

	err := repo.CreateShortURL(rootCtx, mockOriginalURL, mockShortURL)
	require.NoError(t, err, "expected no error when adding short url to repo")

	originalURL, err := service.ResolveShortURL(rootCtx, mockShortURL)
	require.Empty(t, originalURL, "expected empty original url")

	errMsg := "expected invalid short url provided error"
	require.ErrorIs(t, ErrInvalidShortURL, err, errMsg)
}

func TestResolveShortURL_ShortURLNotFound(t *testing.T) {
	rootCtx := context.Background()
	repo := inmemory.New()
	service := New("", repo)

	mockOriginalURL := "http://example.com"
	mockShortURL := ""

	err := repo.CreateShortURL(rootCtx, mockOriginalURL, mockShortURL)
	require.NoError(t, err, "expected no error when adding short url to repo")

	mockShortURL = "Stb1end7"
	originalURL, err := service.ResolveShortURL(rootCtx, mockShortURL)
	require.Empty(t, originalURL, "expected empty original url")

	errMsg := "expected invalid short url provided error"
	require.ErrorIs(t, ErrShortURLNotFound, err, errMsg)
}
