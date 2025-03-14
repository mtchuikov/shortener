package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/mtchuikov/shortener/internal/repo/inmemory"
	"github.com/mtchuikov/shortener/internal/service"
	"github.com/mtchuikov/shortener/pkg/httpheaders"
	"github.com/mtchuikov/shortener/pkg/randtools"
	"github.com/stretchr/testify/require"
)

func TestCreateShortURL_Success(t *testing.T) {
	baseURL := "http://localhost:3214/api/"
	service := service.New(baseURL, inmemory.New())
	handler := New(service)

	rawReqBody := strings.NewReader("http://example.com")
	req := httptest.NewRequest(http.MethodPost, "/", rawReqBody)
	recorder := httptest.NewRecorder()

	handler.CreateShortURL(recorder, req)
	resp := recorder.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 code")
	rawRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "expected no error when reading body")

	url := strings.ReplaceAll(string(rawRespBody), "\n", "")
	var validShortURL = regexp.MustCompile(`^http://localhost:3214/api/[A-Za-z0-9]{8}$`)
	valid := validShortURL.MatchString(url)
	require.True(t, valid, "expected valid short url")
}

func TestCreateShortURL_Success_JSON(t *testing.T) {
	baseURL := "http://localhost:3214/api/"
	service := service.New(baseURL, inmemory.New())
	handler := New(service)

	reqBody := CreateShortURLRequest{
		URL: "http://example.com",
	}
	jsonPayload, err := json.Marshal(reqBody)
	require.NoError(t, err, "expected no error when marshalling body")

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonPayload))
	req.Header.Set(httpheaders.ContentType, httpheaders.ApplicationJSON)
	recorder := httptest.NewRecorder()

	handler.CreateShortURL(recorder, req)
	resp := recorder.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 code")

	var respBody CreateShortURLResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err, "expected no err when decoding response")

	validShortURL := regexp.MustCompile(`^http://localhost:3214/api/[A-Za-z0-9]{8}$`)
	require.True(t, validShortURL.MatchString(respBody.Result), "expected valid short url")
}

func TestCreateShortURL_NoBody(t *testing.T) {
	baseURL := "http://localhost:3214/api/"
	service := service.New(baseURL, inmemory.New())
	handler := New(service)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	recorder := httptest.NewRecorder()

	handler.CreateShortURL(recorder, req)
	resp := recorder.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "expected 400 code for empty body")

	rawRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "expected no error when reading body")

	respBody := strings.ReplaceAll(string(rawRespBody), "\n", "")
	require.Equal(t, ErrNoURLProvided.Error(), respBody, "expected no url provided error")
}

func TestCreateShortURL_TooLongBody(t *testing.T) {
	baseURL := "http://localhost:3214/api/"
	service := service.New(baseURL, inmemory.New())
	handler := New(service)

	path := randtools.DefaultGenerateString(maxURLSize)
	originalURL := "http://example.com/" + path

	rawReqBody := strings.NewReader(originalURL)
	req := httptest.NewRequest(http.MethodPost, "/", rawReqBody)
	recorder := httptest.NewRecorder()

	handler.CreateShortURL(recorder, req)
	resp := recorder.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "expected 400 code for too long body")

	rawRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "expected no error when reading body")

	respBody := strings.ReplaceAll(string(rawRespBody), "\n", "")
	require.Equal(t, ErrURLTooLong.Error(), respBody, "expected url too long error")
}
