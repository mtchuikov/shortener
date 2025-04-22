package handlers

import (
	"net/http"

	"github.com/mtchuikov/shortener/internal/models"
	"github.com/mtchuikov/shortener/internal/repo"

	"github.com/rs/zerolog"
)

func logUnexpectedError(log *zerolog.Logger, err error, op string) {
	log.Error().Err(err).Str("op", op).
		Msg("unexpected error")
}

func modelErrorToCodeAndMsg(err error) (int, string, error) {
	switch err {
	case models.ErrInvalidOriginalURL:
		return http.StatusBadRequest, err.Error(), nil

	case models.ErrInvalidShortenID:
		return http.StatusBadRequest, err.Error(), nil

	default:
		msg := "something went wrong"
		return http.StatusInternalServerError, msg, err
	}
}

func serviceErrorToCodeAndMsg(err error) (int, string, error) {
	switch err {
	case repo.ErrShortenIDAlreadyExists:
		return http.StatusInternalServerError, err.Error(), nil

	case repo.ErrOirignalURLAlreadyShortened:
		return http.StatusConflict, "", nil

	case repo.ErrUnexpectedError:
		return http.StatusInternalServerError, err.Error(), nil

	default:
		msg := "something went wrong"
		return http.StatusInternalServerError, msg, err
	}
}
