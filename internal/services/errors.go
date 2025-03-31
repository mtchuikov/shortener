package services

import "errors"

var (
	ErrInvalidURL  = errors.New("invalid url")
	ErrInvalidID   = errors.New("invalid id")
	ErrIDNotFound  = errors.New("id not found")
	ErrURLNotFound = errors.New("url not found")
)
