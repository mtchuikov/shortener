package jwtauth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type (
	ParseFn     func(context.Context, *JWTAuth) jwt.Keyfunc
	ValidatorFn func(context.Context, *jwt.Token) error
)

type JWTAuth struct {
	Alg          jwt.SigningMethod
	SignKey      any
	VerifyKey    any
	ParseFn      ParseFn
	ValidatorFns []ValidatorFn
}

func New(secretKey any, opts ...Option) *JWTAuth {
	ja := &JWTAuth{
		Alg:          DefaultAlg,
		SignKey:      secretKey,
		VerifyKey:    secretKey,
		ValidatorFns: make([]ValidatorFn, 0),
	}

	ja.ParseFn = DefaultParseFn

	for _, opt := range opts {
		opt(ja)
	}

	return ja
}
