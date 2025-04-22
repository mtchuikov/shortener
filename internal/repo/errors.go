package repo

import "errors"

var (
	ErrShortenIDAlreadyExists      = errors.New("shorten id already exists")
	ErrOirignalURLAlreadyShortened = errors.New("original url already shorten")
	ErrOriginalURLNotFound         = errors.New("original url not found")
	ErrUnexpectedError             = errors.New("unexpected error")
)
