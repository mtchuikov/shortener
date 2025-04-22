package jwtauth

import "errors"

var (
	ErrUnauthorized  = errors.New("token is unauthorized")
	ErrExpired       = errors.New("token is expired")
	ErrNBFInvalid    = errors.New("token nbf validation failed")
	ErrIATInvalid    = errors.New("token iat validation failed")
	ErrNoTokenFound  = errors.New("no token found")
	ErrInvalidAlgo   = errors.New("algorithm mismatch")
	ErrInvalidGoType = errors.New("invalid go type")
)
