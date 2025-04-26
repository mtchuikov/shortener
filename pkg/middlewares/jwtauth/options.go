package jwtauth

import "github.com/golang-jwt/jwt/v5"

type Option func(*JWTAuth)

func WithAlg(alg jwt.SigningMethod) Option {
	return func(a *JWTAuth) {
		a.Alg = alg
	}
}

func WithSignKey(key any) Option {
	return func(a *JWTAuth) {
		a.SignKey = key
	}
}

func WithVerifyKey(key any) Option {
	return func(a *JWTAuth) {
		a.VerifyKey = key
	}
}

func WithValidator(fn ValidatorFn) Option {
	return func(a *JWTAuth) {
		a.ValidatorFns = append(a.ValidatorFns, fn)
	}
}
