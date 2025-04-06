package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mtchuikov/shortener/pkg/pinger"
)

type ping struct {
	pinger *pinger.Pinger
}

func RegisterPing(mux *chi.Mux, pinger *pinger.Pinger) {
	handler := ping{pinger}
	mux.Get("/ping", handler.Handle)
}

func (h *ping) Handle(rw http.ResponseWriter, req *http.Request) {
	err := h.pinger.Error()
	if err != nil {
		errMsg := "failed to ping postgres"
		http.Error(rw, errMsg, http.StatusInternalServerError)
	}
}
