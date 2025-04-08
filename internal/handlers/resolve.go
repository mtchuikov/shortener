package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type resolverService interface {
	Serve(ctx context.Context, shortID string) (string, error)
}

type resolveHandler struct {
	resolver resolverService
}

func RegisterResolve(router chi.Router, service resolverService) {
	handler := resolveHandler{service}
	router.Get("/{short_id}", handler.Handle)
}

func (h *resolveHandler) Handle(rw http.ResponseWriter, req *http.Request) {
	shortID := chi.URLParam(req, "short_id")
	ctx := req.Context()

	originalURL, err := h.resolver.Serve(ctx, shortID)
	if err != nil {
		msg, code := matcErrorToMsgAndCode(err)
		http.Error(rw, msg, code)
		return
	}

	rw.Header().Set("Location", originalURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
