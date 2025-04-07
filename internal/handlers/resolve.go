package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type resolverService interface {
	Serve(ctx context.Context, shortID string) (string, error)
}

func RegisterResolver(router chi.Router, service resolverService) {
	router.Get("/{short_id}", handleResolve(service))
}

func handleResolve(service resolverService) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		shortID := chi.URLParam(req, "short_id")

		originalURL, err := service.Serve(req.Context(), shortID)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		rw.Header().Set("Location", originalURL)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
}
