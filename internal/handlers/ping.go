package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type pingerService interface {
	Error() error
}

func RegisterPinger(router chi.Router, pinger pingerService) {
	router.Get("/ping", handlePing(pinger))
}

func handlePing(service pingerService) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if service.Error() != nil {
			errMsg := "failed to ping postgres"
			http.Error(rw, errMsg, http.StatusInternalServerError)
		}
	}
}
