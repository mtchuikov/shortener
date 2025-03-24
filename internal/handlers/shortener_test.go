package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/mtchuikov/shortener/internal/cache/inmemory"
	"github.com/mtchuikov/shortener/internal/services"
	"github.com/mtchuikov/shortener/pkg/logtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortener_Success(t *testing.T) {
	const (
		mockBaseURL = "http://localhost:3214/short/"
		mockPayload = "https://example.com/"
	)

	logger := logtools.NewZerolog("test", os.Stdout)
	cache := inmemory.New()
	service := services.NewShortener(mockBaseURL, cache)

	handler := shortener{
		logger:  logger,
		service: service,
	}

	body := strings.NewReader(mockPayload)
	req := httptest.NewRequest(http.MethodGet, "/", body)

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	errMsg := "expected status code 201, got %v"
	require.Equalf(t, http.StatusCreated, rr.Code, errMsg, rr.Code)

	payload, err := io.ReadAll(resp.Body)

	errMsg = "expected no error when reading body, got '%v'"
	require.NoErrorf(t, err, errMsg, err)

	rawURL := string(payload)
	url := strings.ReplaceAll(rawURL, "\n", "")

	urlRegexp := regexp.MustCompile(`^` + mockBaseURL + `[A-Za-z0-9]{8}$`)
	valid := urlRegexp.MatchString(url)

	assert.True(t, valid, "expected valid url")
}
