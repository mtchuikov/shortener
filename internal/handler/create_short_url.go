package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mtchuikov/shortener/pkg/httpheaders"
)

type CreateShortURLRequest struct {
	URL string `json:"url"`
}

type CreateShortURLResponse struct {
	Result string `json:"result"`
}

const maxURLSize = 2048

func extractOriginalURL(req *http.Request) (string, bool, error) {
	reader := io.LimitReader(req.Body, maxURLSize+1)

	contentType := req.Header.Get(httpheaders.ContentType)
	isJSON := contentType == httpheaders.ApplicationJSON
	if isJSON {
		var body CreateShortURLRequest
		err := json.NewDecoder(reader).Decode(&body)
		if err != nil {
			return "", isJSON, ErrFailedToReadURLFromBody
		}

		return body.URL, isJSON, nil
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return "", isJSON, ErrFailedToReadURLFromBody
	}

	return string(body), isJSON, nil
}

func initialValidateURL(url string) error {
	urlLen := len(url)
	if urlLen == 0 {
		return ErrNoURLProvided
	}

	if urlLen > maxURLSize {
		return ErrURLTooLong
	}

	return nil
}

func (h *Handler) CreateShortURL(rw http.ResponseWriter, req *http.Request) {
	originalURL, isJSON, err := extractOriginalURL(req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = initialValidateURL(originalURL)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.CreateShortURL(req.Context(), originalURL)
	if err != nil {
		statusCode := serviceErrToStatusCode(err)
		http.Error(rw, err.Error(), statusCode)
		return
	}

	rw.WriteHeader(http.StatusCreated)

	if isJSON {
		rw.Header().Set(httpheaders.ContentType, httpheaders.ApplicationJSON)

		respBody := CreateShortURLResponse{Result: shortURL}
		payload, err := json.Marshal(&respBody)
		if err != nil {
			msg := ""
			http.Error(rw, msg, http.StatusInternalServerError)
			return
		}

		rw.Write(payload)
		return
	}

	rw.Write([]byte(shortURL))
}
