package handler

import "github.com/mtchuikov/shortener/internal/service"

type Handler struct {
	service service.IService
}

func New(service service.IService) *Handler {
	return &Handler{
		service: service,
	}
}
