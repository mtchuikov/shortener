package jwtauth

import "github.com/golang-jwt/jwt/v5"

type JWTAuth struct {
	alg       jwt.SigningMethod
	signKey   any
	verifyKey any
}

func New(alg jwt.SigningMethod, secretKey, verifyKey any) *JWTAuth {
	return &JWTAuth{
		alg:       alg,
		signKey:   secretKey,
		verifyKey: verifyKey,
	}
}
