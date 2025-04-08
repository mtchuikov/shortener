package repo

import "errors"

var (
	ErrFailedToGetOriginalURL = errors.New("failed to get original url")
	ErrOriginalURLNotFound    = errors.New("original url not found")
	ErrShortIDNotFound        = errors.New("short id not found")
)
