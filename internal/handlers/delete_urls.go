package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mtchuikov/shortener/internal/middlewares"
	"github.com/mtchuikov/shortener/internal/services"
)

func RegisterMarkURLsAsDeleted(router chi.Router, srv services.Shortener) {
	router.Delete("/api/user/urls", markURLsAsDeleted(srv))
}

func markURLsAsDeleted(srv services.Shortener) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.Body = http.MaxBytesReader(rw, req.Body, 1088)

		payload, err := io.ReadAll(req.Body)
		if err != nil {
			msg := errMsgPayloadTooLarge
			http.Error(rw, msg, http.StatusBadRequest)
			return
		}

		var slugs []string
		err = json.Unmarshal(payload, &slugs)
		if err != nil {
			msg := errMsgInvalidPayloadFormat
			http.Error(rw, msg, http.StatusBadRequest)
			return
		}

		ctx := req.Context()
		userID, _ := ctx.Value(middlewares.UserIDCtxKey).(string)
		go func() {
			ctx := context.Background()
			timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			srv.MarkShortURLsAsDeleted(timeoutCtx, userID, slugs)
		}()

		rw.WriteHeader(http.StatusAccepted)
	}
}
