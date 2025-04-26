package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mtchuikov/shortener/internal/middlewares"
	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/services"
	"github.com/mtchuikov/shortener/internal/storage"
	"github.com/mtchuikov/shortener/internal/validators"
)

func RegisterShorten(router chi.Router, srv services.Shortener) {
	router.Post("/", shorten(srv))
	router.Post("/api/shorten", shortenJSON(srv))
	router.Post("/api/shorten/batch", shortenBatch(srv))
}

func shorten(srv services.Shortener) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.Body = http.MaxBytesReader(rw, req.Body, 4096)

		payload, err := io.ReadAll(req.Body)
		if err != nil {
			msg := errMsgPayloadTooLarge
			http.Error(rw, msg, http.StatusBadRequest)
			return
		}

		originalURL := string(payload)
		err = validators.OriginalURL(originalURL)
		if err != nil {
			msg := validators.ErrMsgInvalidOriginalURL
			http.Error(rw, msg, http.StatusBadRequest)
			return
		}

		statusCode := http.StatusCreated

		ctx := req.Context()
		userID, _ := ctx.Value(middlewares.UserIDCtxKey).(string)

		shortURL, err := srv.CreateShortURL(ctx, userID, originalURL)
		if err != nil {
			if !errors.Is(err, storage.ErrOriginalURLAlreadyExists) {
				msg := errSomethingWentWrong
				http.Error(rw, msg, http.StatusInternalServerError)
				return
			}

			statusCode = http.StatusConflict
		}

		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(statusCode)
		rw.Write([]byte(shortURL))
	}
}

func shortenJSON(srv services.Shortener) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.Body = http.MaxBytesReader(rw, req.Body, 5006)

		payload, err := io.ReadAll(req.Body)
		if err != nil {
			msg := errMsgPayloadTooLarge
			http.Error(rw, msg, http.StatusBadRequest)
			return
		}

		var reqData models.ShortURL
		err = json.Unmarshal(payload, &reqData)
		if err != nil {
			msg := errMsgInvalidPayloadFormat
			http.Error(rw, msg, http.StatusBadRequest)
			return
		}

		err = validators.OriginalURL(reqData.URL)
		if err != nil {
			msg := validators.ErrMsgInvalidOriginalURL
			http.Error(rw, msg, http.StatusBadRequest)
			return
		}

		statusCode := http.StatusCreated

		ctx := req.Context()
		userID, _ := ctx.Value(middlewares.UserIDCtxKey).(string)

		shortURL, err := srv.CreateShortURL(ctx, userID, reqData.URL)
		if err != nil {
			if !errors.Is(err, storage.ErrOriginalURLAlreadyExists) {
				msg := errSomethingWentWrong
				http.Error(rw, msg, http.StatusInternalServerError)
				return
			}

			statusCode = http.StatusConflict
		}

		respData := models.ShortURLResult{Result: shortURL}
		payload, err = json.Marshal(respData)
		if err != nil {
			msg := errSomethingWentWrong
			http.Error(rw, msg, http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(statusCode)
		rw.Write(payload)
	}
}

func shortenBatch(srv services.Shortener) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.Body = http.MaxBytesReader(rw, req.Body, 131_072)

		payload, err := io.ReadAll(req.Body)
		if err != nil {
			msg := errMsgPayloadTooLarge
			http.Error(rw, msg, http.StatusBadRequest)
			return
		}

		var reqData []models.BatchCreateShortURLs
		err = json.Unmarshal(payload, &reqData)
		if err != nil {
			msg := errMsgInvalidPayloadFormat
			http.Error(rw, msg, http.StatusBadRequest)
			return
		}

		ctx := req.Context()
		userID, _ := ctx.Value(middlewares.UserIDCtxKey).(string)

		result, err := srv.BatchCreateShortURLs(ctx, userID, reqData)
		if err != nil {
			if errors.Is(err, storage.ErrNoItemsInBatch) {
				msg := storage.ErrMsgNoItemsInBatch
				http.Error(rw, msg, http.StatusInternalServerError)
				return
			}

			msg := errSomethingWentWrong
			http.Error(rw, msg, http.StatusInternalServerError)
			return
		}

		payload, err = json.Marshal(result)
		if err != nil {
			msg := errSomethingWentWrong
			http.Error(rw, msg, http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		rw.Write(payload)
	}
}
