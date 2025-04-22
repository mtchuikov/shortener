package handlers

import (
	"context"
	"net/http"

	"github.com/mtchuikov/shortener/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type resolver interface {
	Serve(context.Context, models.ShortenID) (models.OriginalURL, error)
}

type resolve struct {
	log      *zerolog.Logger
	resolver resolver
}

func RegisterResolver(log *zerolog.Logger, router chi.Router, resolver resolver) {
	handler := resolve{
		log:      log,
		resolver: resolver,
	}

	router.Get("/{shorten_id}", handler.Handle)
}

const resolveOp = "handlers.resolve.handle"

func (h *resolve) Handle(rw http.ResponseWriter, req *http.Request) {
	rawShortenID := chi.URLParam(req, "shorten_id")

	shortenID, err := models.NewShortenID(rawShortenID)
	if err != nil {
		code, msg, err := modelErrorToCodeAndMsg(err)
		if err != nil {
			logUnexpectedError(h.log, err, resolveOp)
		}

		http.Error(rw, msg, code)
		return
	}

	originalURL, err := h.resolver.Serve(req.Context(), shortenID)
	if err != nil {
		code, msg, err := serviceErrorToCodeAndMsg(err)
		if err != nil {
			logUnexpectedError(h.log, err, resolveOp)
		}

		http.Error(rw, msg, code)
		return
	}

	rw.Header().Set("Location", originalURL.String())
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
