package handler

import "errors"

var (
	ErrFailedToReadURL = errors.New("failed to read url")
	ErrURLTooLong      = errors.New("url too long")
	ErrNoURLProvided   = errors.New("no url provided")
)
