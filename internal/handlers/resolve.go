package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mtchuikov/shortener/pkg/middlewares"
	"github.com/rs/zerolog"
)

type resolverService interface {
	Serve(ctx context.Context, id string) (string, error)
}

type resolver struct {
	logger  zerolog.Logger
	service resolverService
}

func RegisterResolver(logger zerolog.Logger, mux *chi.Mux, service resolverService) {
	handler := resolver{
		logger:  logger,
		service: service,
	}

	mux.Get("/{short_id}", handler.Handle)
}

func (h *resolver) Handle(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := chi.URLParam(req, "short_id")

	url, err := h.service.Serve(ctx, id)
	if err != nil {
		middlewares.RequestContextWithError(req, err)

		errMsg := errors.Unwrap(err).Error()
		http.Error(rw, errMsg, http.StatusBadRequest)
		return
	}

	rw.Header().Set("Location", url)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
