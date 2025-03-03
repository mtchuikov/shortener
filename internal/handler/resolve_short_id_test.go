package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mtchuikov/shortener/internal/service"
	memstorage "github.com/mtchuikov/shortener/internal/storage/mem"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestResolveShortID_Success(t *testing.T) {
	storage := memstorage.New()

	baseURL := "http://localhost:3214/api/"
	service := service.New(log.Logger, baseURL, storage)
	handler := New(service)

	mux := chi.NewRouter()
	mux.Get("/{short_id}", handler.ResolveShortID)
	mux.Post("/", handler.ShortURL)

	server := httptest.NewServer(mux)
	defer server.Close()
	client := server.Client()

	mockURL := "http://example.com"
	rawReqBody := strings.NewReader(mockURL)
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/", rawReqBody)
	resp, err := client.Do(req)

	require.NoError(t, err, "expected no error when doing request")
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 code")

	rawRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "expected no error when reading body")

	respBody := string(rawRespBody)
	url := strings.ReplaceAll(respBody, "\n", "")
	shortID := url[len(url)-8:]

	req, _ = http.NewRequest(http.MethodGet, server.URL+"/"+shortID, nil)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err = client.Do(req)

	require.NoError(t, err, "expected no error when doing request")
	defer resp.Body.Close()

	require.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode, "expected 307 code")

	location := resp.Header.Get("Location")
	require.Equal(t, mockURL, location, "expected location header to match original url")
}

func TestResolveShortID_InvalidShortID(t *testing.T) {
	storage := memstorage.New()

	baseURL := "http://localhost:3214/api/"
	handler := New(service.New(log.Logger, baseURL, storage))

	mux := chi.NewRouter()
	mux.Get("/{short_id}", handler.ResolveShortID)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := server.Client()
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/not short id", nil)
	resp, err := client.Do(req)

	require.NoError(t, err, "expected no error when doing request")
	defer resp.Body.Close()

	errMsg := "expected 400 status code for invalid short id"
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, errMsg)

	rawRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "expected no error when reading body")

	respBody := strings.ReplaceAll(string(rawRespBody), "\n", "")
	errMsg = "expected no short id provided error"
	require.Equal(t, service.ErrInvalidID.Error(), respBody, errMsg)
}

func TestResolveShortID_NoSuchShortID(t *testing.T) {
	storage := memstorage.New()

	baseURL := "http://localhost:3214/api/"
	handler := New(service.New(log.Logger, baseURL, storage))

	mux := chi.NewRouter()
	mux.Get("/{short_id}", handler.ResolveShortID)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := server.Client()
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/dtb1end7", nil)
	resp, err := client.Do(req)

	require.NoError(t, err, "expected no error when doing request")
	defer resp.Body.Close()

	errMsg := "expected 404 status code for non-existing short id"
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, errMsg)

	rawRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "expected no error when reading body")

	respBody := strings.ReplaceAll(string(rawRespBody), "\n", "")
	errMsg = "expected no such short id error"
	require.Equal(t, service.ErrIDNotFound.Error(), respBody, errMsg)
}
