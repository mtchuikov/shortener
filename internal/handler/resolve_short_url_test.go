package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mtchuikov/shortener/internal/repo/inmemory"
	"github.com/mtchuikov/shortener/internal/service"
	"github.com/stretchr/testify/require"
)

func TestResolveShortURL_Success(t *testing.T) {
	baseURL := "http://localhost:3214/api/"
	service := service.New(baseURL, inmemory.New())
	handler := New(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.CreateShortURL)
	mux.HandleFunc("/{id}", handler.ResolveShortURL)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := server.Client()

	originalURL := "http://example.com"
	rawReqBody := strings.NewReader(originalURL)
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/", rawReqBody)
	resp, err := client.Do(req)

	require.NoError(t, err, "expected no error when doing request")
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 code")

	rawRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "expected no error when reading body")

	respBody := string(rawRespBody)
	url := strings.ReplaceAll(respBody, "\n", "")
	shortURL := url[len(url)-8:]

	req, _ = http.NewRequest(http.MethodGet, server.URL+"/"+shortURL, nil)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err = client.Do(req)

	require.NoError(t, err, "expected no error when doing request")
	defer resp.Body.Close()

	require.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode, "expected 307 code")

	location := resp.Header.Get("Location")
	require.Equal(t, originalURL, location, "expected location header to match original url")
}

func TestResolveShortURL_InvalidShortURL(t *testing.T) {
	handler := New(service.New("", inmemory.New()))

	mux := http.NewServeMux()
	mux.HandleFunc("/{id}", handler.ResolveShortURL)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := server.Client()
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/invalid", nil)
	resp, err := client.Do(req)

	require.NoError(t, err, "expected no error when doing request")
	defer resp.Body.Close()

	errMsg := "expected 400 status code for invalid short id"
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, errMsg)

	rawRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "expected no error when reading body")

	respBody := strings.ReplaceAll(string(rawRespBody), "\n", "")
	errMsg = "expected invalid short url provided error"
	require.Equal(t, service.ErrInvalidShortURL.Error(), respBody, errMsg)
}

func TestResolveShortURL_ShortURLNotFound(t *testing.T) {
	baseURL := "http://localhost:3214/api/"
	handler := New(service.New(baseURL, inmemory.New()))

	mux := http.NewServeMux()
	mux.HandleFunc("/{id}", handler.ResolveShortURL)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := server.Client()
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/Stb1end7", nil)
	resp, err := client.Do(req)

	require.NoError(t, err, "expected no error when doing request")
	defer resp.Body.Close()

	errMsg := "expected 400 status code for non-existing short id"
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, errMsg)

	rawRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "expected no error when reading body")

	respBody := strings.ReplaceAll(string(rawRespBody), "\n", "")
	require.Equal(t, service.ErrShortURLNotFound.Error(), respBody, "expected short url not found error")
}
