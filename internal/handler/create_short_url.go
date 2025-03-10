package handler

import (
	"io"
	"net/http"
)

const maxURLSize = 2048

func (h *Handler) CreateShortURL(rw http.ResponseWriter, req *http.Request) {
	reader := io.LimitReader(req.Body, maxURLSize+1)
	body, err := io.ReadAll(reader)
	if err != nil {
		msg := ErrFailedToReadURLFromBody.Error()
		http.Error(rw, msg, http.StatusBadRequest)
		return
	}

	bodyLen := len(body)
	if bodyLen == 0 {
		msg := ErrNoURLProvided.Error()
		http.Error(rw, msg, http.StatusBadRequest)
		return
	}

	if bodyLen > maxURLSize {
		msg := ErrURLTooLong.Error()
		http.Error(rw, msg, http.StatusBadRequest)
		return
	}

	originalURL := string(body)
	shortURL, err := h.service.CreateShortURL(req.Context(), originalURL)
	if err != nil {
		statusCode := serviceErrToStatusCode(err)
		http.Error(rw, err.Error(), statusCode)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(shortURL))
}
