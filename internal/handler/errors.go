package handler

import (
	"errors"
	"net/http"

	"github.com/mtchuikov/shortener/internal/service"
)

var (
	ErrFailedToReadBody = errors.New("failed to read body")
	ErrURLTooLong       = errors.New("url too long")
	ErrNoURLProvided    = errors.New("no url provided")
)

func serviceErrToStatusCode(err error) int {
	if errors.Is(err, service.ErrInternalError) {
		return http.StatusInternalServerError
	}

	return http.StatusBadRequest
}
