package jwtauth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

var (
	DefaultAlg     = jwt.SigningMethodHS256
	DefaultParseFn = func(ctx context.Context, ja *JWTAuth) jwt.Keyfunc {
		return func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != ja.Alg.Alg() {
				return nil, ErrInvalidAlgo
			}

			return ja.VerifyKey, nil
		}
	}
)
