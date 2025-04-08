package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type pingerService interface {
	Serve(ctx context.Context) error
}

type pingHandler struct {
	pinger pingerService
}

func RegisterPing(router chi.Router, service pingerService) {
	handler := pingHandler{service}
	router.Get("/ping", handler.Handle)
}

func (h *pingHandler) Handle(rw http.ResponseWriter, req *http.Request) {
	err := h.pinger.Serve(req.Context())
	if err != nil {
		msg := "failed to ping"
		http.Error(rw, msg, http.StatusInternalServerError)
	}
}
