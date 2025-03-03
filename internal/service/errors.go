package service

import "errors"

var (
	ErrInvalidURL    = errors.New("invalid url")
	ErrInvalidID     = errors.New("invalid id")
	ErrIDNotFound    = errors.New("id not found")
	ErrInternalError = errors.New("internal error")
)
