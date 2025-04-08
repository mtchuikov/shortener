package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mtchuikov/shortener/internal/models"
)

type shortenerService interface {
	Serve(ctx context.Context, originalURL string) (string, error)
	ServeBatch(
		ctx context.Context,
		urlsToShort models.URLsToShort,
		numUrlsToShort int,
	) (models.ShortenURLs, error)
}

type shortHandler struct {
	urlMaxLen int64
	shortener shortenerService
}

func RegisterShortener(router chi.Router, service shortenerService) {
	handler := shortHandler{
		urlMaxLen: 4096 + 1,
		shortener: service,
	}

	router.Post("/", handler.Handle)
	router.Post("/api/shorten", handler.Handle)
	router.Post("/api/shorten/batch", handler.HandleBatch)
}

func (h *shortHandler) extractURL(body io.Reader, isJSON bool) (string, string) {
	// for plain text requests, the url can occupy the full 2048 bytes
	// since the body contains only the url, but for json requests, the
	// url must be smaller cause json structure includes extra chars
	// like "{}" etc
	limitReader := io.LimitReader(body, h.urlMaxLen)

	payload, err := io.ReadAll(limitReader)
	if err != nil {
		return "", "failed to read body"
	}

	if !isJSON {
		url := string(payload)
		return url, ""
	}

	var data models.URLToShort
	err = json.Unmarshal(payload, &data)
	if err != nil {
		return "", "failed to unmarshal"
	}

	return data.URL, ""
}

func (h *shortHandler) Handle(rw http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	isJSON := contentType == "application/json"

	originalURL, msg := h.extractURL(req.Body, isJSON)
	if msg != "" {
		http.Error(rw, msg, http.StatusBadRequest)
		return
	}

	code := http.StatusCreated
	ctx := req.Context()

	shortURL, err := h.shortener.Serve(ctx, originalURL)
	if err != nil {
		msg, code = matcErrorToMsgAndCode(err)
		if code != http.StatusConflict {
			http.Error(rw, msg, code)
			return
		}
	}

	if !isJSON {
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(code)
		rw.Write([]byte(shortURL))
		return
	}

	data := models.ShortenURL{Result: shortURL}
	payload, err := json.Marshal(data)
	if err != nil {
		msg = "failed to marshal"
		http.Error(rw, msg, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	rw.Write(payload)
}

func (h *shortHandler) extractBatchURL(body io.Reader) (models.URLsToShort, string) {
	payload, err := io.ReadAll(body)
	if err != nil {
		return nil, "failed to read body"
	}

	var data models.URLsToShort
	err = json.Unmarshal(payload, &data)
	if err != nil {
		return nil, "failed to unmarshal"
	}

	return data, ""
}

func (h *shortHandler) HandleBatch(rw http.ResponseWriter, req *http.Request) {
	urlsToShort, msg := h.extractBatchURL(req.Body)
	if msg != "" {
		http.Error(rw, msg, http.StatusBadRequest)
		return
	}

	numURLsToShort := len(urlsToShort)
	if numURLsToShort == 0 {
		msg = "no items in batch"
		http.Error(rw, msg, http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	shortenURLs, err := h.shortener.ServeBatch(ctx, urlsToShort, numURLsToShort)
	if err != nil {
		msg, code := matcErrorToMsgAndCode(err)
		http.Error(rw, msg, code)
		return
	}

	payload, err := json.Marshal(shortenURLs)
	if err != nil {
		msg = "failed to marshal"
		http.Error(rw, msg, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(payload)
}
