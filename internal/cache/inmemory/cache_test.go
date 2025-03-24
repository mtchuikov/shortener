package inmemory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	const (
		mockURL = "https://example.com"
		mockID  = "abc123"
	)

	cache := New()

	ctx := context.Background()
	err := cache.CreateShortURL(ctx, mockURL, mockID)

	errMsg := "expected no error when creating short url, got %v"
	require.NoErrorf(t, err, errMsg, err)

	url, err := cache.GetURL(ctx, mockID)

	errMsg = "expected no error when getting url, got %v"
	require.NoErrorf(t, err, errMsg, err)

	errMsg = "expected url and mock url to be equal"
	assert.Equal(t, mockURL, url, errMsg)

	id, err := cache.GetID(ctx, mockURL)

	errMsg = "expected no error when getting id, got %v"
	require.NoErrorf(t, err, errMsg, err)

	errMsg = "expected id and mock id to be equal"
	assert.Equal(t, mockID, id, errMsg)
}
