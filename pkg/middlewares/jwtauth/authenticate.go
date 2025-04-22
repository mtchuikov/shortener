package jwtauth

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func matchParseError(err error) error {
	if errors.Is(err, jwt.ErrTokenExpired) {
		return ErrExpired
	}

	if errors.Is(err, jwt.ErrTokenNotValidYet) {
		return ErrNBFInvalid
	}

	return ErrUnauthorized
}

func Authenticate(ja *JWTAuth, extractor TokenExtractor) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(rw http.ResponseWriter, req *http.Request) {
			tokenString := extractor(req)
			if tokenString == "" {
				msg := ErrNoTokenFound.Error()
				http.Error(rw, msg, http.StatusUnauthorized)
				return
			}

			parseFn := func(t *jwt.Token) (any, error) {
				if t.Method.Alg() != ja.alg.Alg() {
					return nil, ErrInvalidAlgo
				}

				return ja.verifyKey, nil
			}

			token, err := jwt.Parse(tokenString, parseFn)
			if err != nil {
				err = matchParseError(err)
				http.Error(rw, err.Error(), http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				msg := ErrUnauthorized.Error()
				http.Error(rw, msg, http.StatusUnauthorized)
				return
			}

			ctx := PassTokenToContext(req.Context(), token, nil)
			req = req.WithContext(ctx)

			next.ServeHTTP(rw, req)
		}

		return http.HandlerFunc(hfn)
	}
}
