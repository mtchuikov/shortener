package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/mtchuikov/shortener/internal/service"
	memstorage "github.com/mtchuikov/shortener/internal/storage/mem"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestCreateShortURL_Success(t *testing.T) {
	storage := memstorage.New()

	baseURL := "http://localhost:3214/api/"
	service := service.New(log.Logger, baseURL, storage)
	handler := New(service)

	rawReqBody := strings.NewReader("http://example.com")
	req := httptest.NewRequest(http.MethodPost, "/", rawReqBody)
	recorder := httptest.NewRecorder()

	handler.ShortURL(recorder, req)
	resp := recorder.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 code")
	rawRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "expected no error when reading body")

	url := strings.ReplaceAll(string(rawRespBody), "\n", "")
	var validShortURL = regexp.MustCompile(`^http://localhost:3214/api/[A-Za-z0-9]{8}$`)
	valid := validShortURL.MatchString(url)
	fmt.Println(url)
	require.True(t, valid, "expected valid short url")
}

func TestCreateShortURL_NoBody(t *testing.T) {
	storage := memstorage.New()

	baseURL := "http://localhost:3214/api/"
	service := service.New(log.Logger, baseURL, storage)
	handler := New(service)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	recorder := httptest.NewRecorder()

	handler.ShortURL(recorder, req)
	resp := recorder.Result()
	defer resp.Body.Close()

	errMsg := "expected 400 code for empty body"
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, errMsg)

	rawRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "expected no error when reading body")

	respBody := strings.ReplaceAll(string(rawRespBody), "\n", "")
	errMsg = "expected no url provided error"
	require.Equal(t, ErrNoURLProvided.Error(), respBody, errMsg)
}
