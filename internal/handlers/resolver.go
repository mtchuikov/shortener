package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type resolverService interface {
	Serve(ctx context.Context, id string) (string, error)
}

type resolver struct {
	logger  zerolog.Logger
	service resolverService
}

func RegisterResolver(lg zerolog.Logger, mux *chi.Mux, srv resolverService) {
	handler := resolver{
		logger:  lg,
		service: srv,
	}

	mux.Get("/{short_id}", handler.Handle)
}

func (h *resolver) Handle(rw http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "short_id")
	url, err := h.service.Serve(req.Context(), id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Location", url)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
