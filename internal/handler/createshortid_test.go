package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/mtchuikov/shortener/internal/service"
	storage "github.com/mtchuikov/shortener/internal/storage/inmemory"
	"github.com/mtchuikov/shortener/pkg/randtools"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestCreateShortID_Success(t *testing.T) {
	storage := storage.New()
	service := service.New(log.Logger, storage)
	handler := New(service)

	rawReqBody := strings.NewReader("http://example.com")
	req := httptest.NewRequest(http.MethodPost, "/", rawReqBody)
	recorder := httptest.NewRecorder()

	handler.CreateShortID(recorder, req)
	resp := recorder.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 code")

	rawRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "expected no error when reading body")

	respBody := strings.ReplaceAll(string(rawRespBody), "\n", "")
	validURLRegex := regexp.MustCompile(`^(https?://)[^/]+/[A-Za-z0-9]{8}$`)
	require.True(t, validURLRegex.MatchString(respBody), "expected valid url, got %q", respBody)
}

func TestCreateShortID_NoBody(t *testing.T) {
	storage := storage.New()
	service := service.New(log.Logger, storage)
	handler := New(service)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	recorder := httptest.NewRecorder()

	handler.CreateShortID(recorder, req)
	resp := recorder.Result()

	require.Equal(t, http.StatusNotFound, resp.StatusCode, "expected 404 code for empty body")

	rawRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "expected no error when reading body")

	respBody := strings.ReplaceAll(string(rawRespBody), "\n", "")
	require.Equal(t, ErrNoURLProvided.Error(), respBody, "expected no url provided error")
}

func TestCreateShortID_TooLongBody(t *testing.T) {
	storage := storage.New()
	service := service.New(log.Logger, storage)
	handler := New(service)

	path := randtools.DefaultGenerateString(maxURLSize)
	originalURL := "http://example.com/" + path

	rawReqBody := strings.NewReader(originalURL)
	req := httptest.NewRequest(http.MethodPost, "/", rawReqBody)
	recorder := httptest.NewRecorder()

	handler.CreateShortID(recorder, req)
	resp := recorder.Result()

	errMsg := "expected 413 code for too long body"
	require.Equal(t, http.StatusRequestEntityTooLarge, resp.StatusCode, errMsg)

	rawRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "expected no error when reading body")

	respBody := strings.ReplaceAll(string(rawRespBody), "\n", "")
	require.Equal(t, ErrURLTooLong.Error(), respBody, "expected url too long error")
}
