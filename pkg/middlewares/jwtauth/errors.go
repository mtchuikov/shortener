package jwtauth

import "errors"

const (
	ErrMsgInvalidAlgo   = "algorithm mismatch"
	ErrMsgInvalidGoType = "invalid go type"
	ErrMsgNoTokenFound  = "no token found"
	ErrMsgExpired       = "token expired"
	ErrMsgNBFInvalid    = "token nbf validation failed"
	ErrMsgIATInvalid    = "token iat validation failed"
	ErrMsgUnauthorized  = "token unauthorized"
)

var (
	ErrInvalidAlgo   = errors.New(ErrMsgInvalidAlgo)
	ErrInvalidGoType = errors.New(ErrMsgInvalidGoType)
)
