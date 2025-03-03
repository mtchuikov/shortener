package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ResolveShortID(rw http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "short_id")
	url, err := h.service.ResolveShortID(req.Context(), id)
	if err != nil {
		code := serviceErrToStatusCode(err)
		http.Error(rw, err.Error(), code)
		return
	}

	rw.Header().Set("Location", url)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
