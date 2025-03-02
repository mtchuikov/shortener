package handler

import (
	"io"
	"net/http"
)

const maxURLSize = 2048 + 1

func (h *Handler) CreateShortID(rw http.ResponseWriter, req *http.Request) {
	reader := io.LimitReader(req.Body, maxURLSize)
	body, err := io.ReadAll(reader)
	if err != nil {
		errMsg := ErrFailedToReadURL.Error()
		http.Error(rw, errMsg, http.StatusBadRequest)
		return
	}

	bodyLen := len(body)
	if bodyLen == 0 {
		errMsg := ErrNoURLProvided.Error()
		http.Error(rw, errMsg, http.StatusBadRequest)
		return
	}

	if bodyLen >= maxURLSize {
		errMsg := ErrURLTooLong.Error()
		http.Error(rw, errMsg, http.StatusRequestEntityTooLarge)
		return
	}

	ctx := req.Context()
	originalURL := string(body)
	shortID, code, err := h.service.CreateShortID(ctx, originalURL)
	if err != nil {
		http.Error(rw, err.Error(), code)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte("http://localhost:8080/" + shortID))
}
