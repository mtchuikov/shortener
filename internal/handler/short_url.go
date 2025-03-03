package handler

import (
	"io"
	"net/http"
)

const maxURLSize = 2048

func (h *Handler) ShortURL(rw http.ResponseWriter, req *http.Request) {
	reader := io.LimitReader(req.Body, maxURLSize+1)
	rawBody, err := io.ReadAll(reader)
	if err != nil {
		msg := ErrFailedToReadBody.Error()
		http.Error(rw, msg, http.StatusBadRequest)
		return
	}

	rawBodyLen := len(rawBody)
	if rawBodyLen == 0 {
		msg := ErrNoURLProvided.Error()
		http.Error(rw, msg, http.StatusBadRequest)
		return
	}

	if rawBodyLen >= maxURLSize {
		msg := ErrURLTooLong.Error()
		http.Error(rw, msg, http.StatusBadRequest)
		return
	}

	url := string(rawBody)
	shortURL, err := h.service.ShortURL(req.Context(), url)
	if err != nil {
		code := serviceErrToStatusCode(err)
		http.Error(rw, err.Error(), code)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(shortURL))
}
