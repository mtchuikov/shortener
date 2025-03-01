package service

import "errors"

var (
	ErrInvalidURL        = errors.New("invalid url")
	ErrInvalidShortID    = errors.New("invalid short id")
	ErrNoSuchShortID     = errors.New("no such short id")
	ErrSomethidWentWrong = errors.New("something went wrong")
)
