package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ResolveShortURL(rw http.ResponseWriter, req *http.Request) {
	shortURL := chi.URLParam(req, "id")
	originalURL, err := h.service.ResolveShortURL(req.Context(), shortURL)
	if err != nil {
		statusCode := serviceErrToStatusCode(err)
		http.Error(rw, err.Error(), statusCode)
		return
	}

	rw.Header().Set("Location", originalURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
