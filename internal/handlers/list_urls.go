package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mtchuikov/shortener/internal/middlewares"
	"github.com/mtchuikov/shortener/internal/services"
	"github.com/mtchuikov/shortener/internal/storage"
)

func RegisterListURLs(router chi.Router, srv services.Shortener) {
	router.Get("/api/user/urls", listURLs(srv))
}

func listURLs(srv services.Shortener) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		userID, _ := ctx.Value(middlewares.UserIDCtxKey).(string)

		result, err := srv.ListShortURLsByUser(ctx, userID)
		if err != nil {
			if errors.Is(err, storage.ErrUserHasNoShortURLs) {
				msg := storage.ErrMsgUserHasNoShortURLs
				http.Error(rw, msg, http.StatusNoContent)
				return
			}

			msg := errSomethingWentWrong
			http.Error(rw, msg, http.StatusInternalServerError)
			return
		}

		payload, err := json.Marshal(result)
		if err != nil {
			msg := errSomethingWentWrong
			http.Error(rw, msg, http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(payload)
	}
}
