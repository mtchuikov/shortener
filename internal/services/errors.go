package services

import "errors"

var (
	ErrInvalidOriginalURL        = errors.New("invalid original url")
	ErrOriginalURLAlreadyShorten = errors.New("original url already shorten")
	ErrInvalidShortID            = errors.New("invalid short id")
)
