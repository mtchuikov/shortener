package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type pinger interface {
	Serve(context.Context) error
}

type ping struct {
	log    *zerolog.Logger
	pinger pinger
}

func RegisterPing(log *zerolog.Logger, router chi.Router, pinger pinger) {
	handler := ping{
		log:    log,
		pinger: pinger,
	}

	router.Get("/ping", handler.Handle)
}

func (h *ping) Handle(rw http.ResponseWriter, req *http.Request) {
	err := h.pinger.Serve(req.Context())
	if err != nil {
		code, msg, err := serviceErrorToCodeAndMsg(err)
		if err != nil {
			logUnexpectedError(h.log, err, resolveOp)
		}

		http.Error(rw, msg, code)
		return
	}

	rw.Write([]byte("pong"))
}
