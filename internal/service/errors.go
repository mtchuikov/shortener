package service

import "errors"

var (
	ErrInternalError      = errors.New("internal error")
	ErrInvalidOriginalURL = errors.New("invalid  original url")
	ErrInvalidShortURL    = errors.New("invalid short url")
	ErrShortURLNotFound   = errors.New("short url not found")
)
