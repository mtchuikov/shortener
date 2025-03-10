package handler

import "net/http"

func (h *Handler) ResolveShortURL(rw http.ResponseWriter, req *http.Request) {
	shortURL := req.PathValue("id")
	originalURL, err := h.service.ResolveShortURL(req.Context(), shortURL)
	if err != nil {
		statusCode := serviceErrToStatusCode(err)
		http.Error(rw, err.Error(), statusCode)
		return
	}

	rw.Header().Set("Location", originalURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
