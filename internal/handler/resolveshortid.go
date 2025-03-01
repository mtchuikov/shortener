package handler

import (
	"net/http"
)

func (h *Handler) ResolveShortID(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	shortID := req.PathValue("short_id")
	originalURL, code, err := h.service.ResolveShortID(ctx, shortID)
	if err != nil {
		http.Error(rw, err.Error(), code)
		return
	}

	code = http.StatusTemporaryRedirect
	rw.Header().Set("Location", originalURL)
	http.Redirect(rw, req, originalURL, code)
}
