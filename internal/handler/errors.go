package handler

import (
	"errors"
	"net/http"

	"github.com/mtchuikov/shortener/internal/service"
)

var (
	ErrFailedToReadURLFromBody = errors.New("failed to read url from body")
	ErrURLTooLong              = errors.New("url too long")
	ErrNoURLProvided           = errors.New("no url provided")
)

func serviceErrToStatusCode(err error) int {
	if errors.Is(err, service.ErrInternalError) {
		return http.StatusInternalServerError
	}

	return http.StatusBadRequest
}
