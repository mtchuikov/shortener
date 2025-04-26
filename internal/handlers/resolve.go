package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mtchuikov/shortener/internal/services"
	"github.com/mtchuikov/shortener/internal/storage"
	"github.com/mtchuikov/shortener/internal/validators"
)

func RegisterResolve(router chi.Router, srv services.Shortener) {
	router.Get("/{slug}", resolve(srv))
}

func resolve(srv services.Shortener) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		slug := chi.URLParam(req, "slug")

		err := validators.Slug(slug)
		if err != nil {
			msg := validators.ErrMsgInvalidSlug
			http.Error(rw, msg, http.StatusBadRequest)
			return
		}

		originalURL, deleted, err := srv.GetOriginalURLBySlug(req.Context(), slug)
		if err != nil {
			if err == storage.ErrOriginalURLNotFound {
				msg := storage.ErrMsgOriginalURLNotFound
				http.Error(rw, msg, http.StatusBadRequest)
				return
			}

			msg := errSomethingWentWrong
			http.Error(rw, msg, http.StatusInternalServerError)
			return
		}

		if deleted {
			rw.WriteHeader(http.StatusGone)
			return
		}

		rw.Header().Set("Location", originalURL)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
}
