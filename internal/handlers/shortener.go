package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

type shortenerService interface {
	Serve(ctx context.Context, url string) (string, error)
}

type shortener struct {
	logger  zerolog.Logger
	service shortenerService
}

func RegisterShortener(lg zerolog.Logger, mux *chi.Mux, srv shortenerService) {
	handler := shortener{
		logger:  lg,
		service: srv,
	}

	mux.Post("/", handler.Handle)
	mux.Post("/api/shorten", handler.Handle)
}

var (
	ErrFailedToReadBody      = errors.New("failed to read body")
	ErrFailedToUnmarshalJSON = errors.New("failed to unmarshal json")
)

func (h *shortener) extractURL(body io.Reader, isJSON bool) (string, error) {
	const limit = 2048 + 1
	// for plain text requests, the url can occupy the full 2048 bytes
	// since the body contains only the url, but for json requests, the
	// url must be smaller cause json structure includes extra chars
	// like "{}" etc
	lr := io.LimitReader(body, limit)

	payload, err := io.ReadAll(lr)
	if err != nil {
		return "", ErrFailedToReadBody
	}

	if isJSON {
		var data shortenerRequest
		err = json.Unmarshal(payload, &data)
		if err != nil {
			return "", ErrFailedToUnmarshalJSON
		}

		return data.URL, nil
	}

	url := string(payload)
	return url, nil
}

type shortenerRequest struct {
	URL string `json:"url"`
}

type shortenerResponse struct {
	Result string `json:"result"`
}

func (h *shortener) Handle(rw http.ResponseWriter, req *http.Request) {
	ct := req.Header.Get("Content-Type")
	isJSON := ct == "application/json"

	url, err := h.extractURL(req.Body, isJSON)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.Serve(req.Context(), url)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if isJSON {
		json := shortenerResponse{Result: shortURL}
		payload, err := jsoniter.Marshal(&json)
		if err != nil {
			errMsg := "failed to marshal json"
			http.Error(rw, errMsg, http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		rw.Write(payload)
		return
	}

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(shortURL))
}
