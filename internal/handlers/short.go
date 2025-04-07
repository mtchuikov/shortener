package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type shortenerService interface {
	Serve(ctx context.Context, originalURL string) (string, error)
}

func RegisterShortener(
	log zerolog.Logger,
	router chi.Router,
	service shortenerService,
) {
	router.Post("/", handleShort(service))
	router.Post("/api/shorten", nil)
}

const maxURLLen = 4096 + 1

func readBody(body io.Reader) ([]byte, error) {
	// for plain text requests, the url can occupy the full 2048 bytes
	// since the body contains only the url, but for json requests, the
	// url must be smaller cause json structure includes extra chars
	// like "{}" etc
	limitReader := io.LimitReader(body, maxURLLen)
	return io.ReadAll(limitReader)
}

type shortenerRequest struct {
	URL string `json:"url"`
}

const wFailedReadRequestBody = "failed to read request body: %w"
const wFailedToUnmarshalJSON = "failed to unmarshal json: %w"

func extractURL(body io.Reader, isJSON bool) (string, error) {
	payload, err := readBody(body)
	if err != nil {
		return "", fmt.Errorf(wFailedReadRequestBody, err)
	}

	if isJSON {
		var data shortenerRequest
		err = json.Unmarshal(payload, &data)
		if err != nil {
			return "", fmt.Errorf(wFailedToUnmarshalJSON, err)
		}

		return data.URL, nil
	}

	url := string(payload)
	return url, nil
}

type shortenerResponse struct {
	Result string `json:"result"`
}

func handleShort(service shortenerService) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ct := req.Header.Get("Content-Type")
		isJSON := ct == "application/json"

		originalURL, err := extractURL(req.Body, isJSON)
		if err != nil {
			errMsg := "failed to read body"
			http.Error(rw, errMsg, http.StatusBadRequest)
			return
		}

		shortURL, err := service.Serve(req.Context(), originalURL)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		if isJSON {
			data := shortenerResponse{Result: shortURL}
			payload, err := json.Marshal(&data)
			if err != nil {
				http.Error(rw, "", http.StatusInternalServerError)
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
}
