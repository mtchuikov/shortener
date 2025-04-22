package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/mtchuikov/shortener/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type shortener interface {
	Serve(context.Context, models.OriginalURL) (models.ShortenURL, error)
	ServeBatch(
		ctx context.Context, urlsToShort models.BatchShortURLs,
		batchSize int,
	) (
		models.BatchShortenURLs,
		error,
	)
}

type shorten struct {
	log            *zerolog.Logger
	shortener      shortener
	maxURLSize     int64
	maxURLSizeJSON int64
}

func RegisterShorten(log *zerolog.Logger, router chi.Router, shortener shortener) {
	handler := shorten{
		log:            log,
		shortener:      shortener,
		maxURLSize:     4096,
		maxURLSizeJSON: 4096 + 10,
	}

	router.Post("/", handler.Handle)
	router.Post("/api/shorten", handler.HandleJSON)
	router.Post("/api/shorten/batch", handler.HandleJSONBatch)
}

func (h *shorten) Handle(rw http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(rw, req.Body, h.maxURLSize)
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "payload too large", http.StatusBadRequest)
		return
	}

	url := string(payload)
	originalURL, err := models.NewOriginalURL(url)
	if err != nil {
		code, msg, err := modelErrorToCodeAndMsg(err)
		if err != nil {
			logUnexpectedError(h.log, err, resolveOp)
		}

		http.Error(rw, msg, code)
		return
	}

	statusCode := http.StatusCreated

	shortenURL, err := h.shortener.Serve(req.Context(), originalURL)
	if err != nil {
		code, msg, err := serviceErrorToCodeAndMsg(err)
		if err != nil {
			logUnexpectedError(h.log, err, resolveOp)
		}

		if code != http.StatusConflict {
			http.Error(rw, msg, code)
			return
		}

		statusCode = http.StatusConflict
	}

	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(statusCode)
	rw.Write([]byte(shortenURL))
}
