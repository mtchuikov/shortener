package handlers

import (
	"net/http"

	"github.com/mtchuikov/shortener/internal/repo"
	"github.com/mtchuikov/shortener/internal/services"
)

func matcErrorToMsgAndCode(err error) (string, int) {
	var msg string
	switch err {
	case services.ErrInvalidOriginalURL:
		msg = "invalid original url"
		return msg, http.StatusBadRequest

	case services.ErrOriginalURLAlreadyShorten:
		return "", http.StatusConflict

	case services.ErrInvalidShortID:
		msg = "invalid short id"
		return msg, http.StatusBadRequest

	case repo.ErrFailedToGetOriginalURL:
		msg = "failed to get original url"
		return msg, http.StatusInternalServerError

	case repo.ErrOriginalURLNotFound:
		msg = "original url not found"
		return msg, http.StatusBadRequest

	case repo.ErrShortIDNotFound:
		msg = "short id not found"
		return msg, http.StatusBadRequest

	default:
		msg = "something went wrong"
		return msg, http.StatusInternalServerError
	}
}
