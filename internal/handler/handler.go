package handler

import (
	"github.com/mtchuikov/shortener/internal/interfaces"
)

type Handler struct {
	service interfaces.Service
}

func New(service interfaces.Service) *Handler {
	return &Handler{
		service: service,
	}
}
