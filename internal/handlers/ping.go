package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type pinger interface {
	Ping(context.Context) error
}

func RegisterPing(router chi.Router, pinger pinger) {
	router.Get("/ping", ping(pinger))
}

func ping(srv pinger) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		err := srv.Ping(req.Context())
		if err != nil {
			msg := errSomethingWentWrong
			http.Error(rw, msg, http.StatusInternalServerError)
			return
		}

		rw.Write([]byte("pong"))
	}
}
