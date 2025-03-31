package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mtchuikov/shortener/pkg/middlewares"
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

func (h *shortener) extractURL(body io.Reader, isJSON bool) (string, error) {
	const op = "handler.shortener.extract_url"

	const limit = 2048 + 1
	// for plain text requests, the url can occupy the full 2048 bytes
	// since the body contains only the url, but for json requests, the
	// url must be smaller cause json structure includes extra chars
	// like "{}" etc
	lr := io.LimitReader(body, limit)

	payload, err := io.ReadAll(lr)
	if err != nil {
		return "", fmt.Errorf(
			"%s - %w: %s",
			op, errFailedToReadBody, err,
		)
	}

	if isJSON {
		var data shortenerRequest
		err = json.Unmarshal(payload, &data)
		if err != nil {
			return "", fmt.Errorf(
				"%s - %w: %s",
				op, errFailedToUnmarshalJSON, err,
			)
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
	const op = "handler.shortener.handle"

	ct := req.Header.Get("Content-Type")
	isJSON := ct == "application/json"

	ctx := req.Context()

	url, err := h.extractURL(req.Body, isJSON)
	if err != nil {
		middlewares.RequestContextWithError(req, err)

		errMsg := errors.Unwrap(err).Error()
		http.Error(rw, errMsg, http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.Serve(ctx, url)
	if err != nil {
		middlewares.RequestContextWithError(req, err)

		errMsg := errors.Unwrap(err).Error()
		http.Error(rw, errMsg, http.StatusBadRequest)
		return
	}

	if isJSON {
		data := shortenerResponse{Result: shortURL}
		payload, err := json.Marshal(&data)
		if err != nil {
			err = fmt.Errorf(
				"%s - %w: %s",
				op, errFailedToMarshalJSON, err,
			)
			middlewares.RequestContextWithError(req, err)

			errMsg := errFailedToMarshalJSON.Error()
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
