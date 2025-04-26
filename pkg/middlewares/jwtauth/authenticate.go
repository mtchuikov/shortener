package jwtauth

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func parseErrorToMsg(err error) string {
	if errors.Is(err, jwt.ErrTokenExpired) {
		return ErrMsgExpired
	}

	if errors.Is(err, jwt.ErrTokenNotValidYet) {
		return ErrMsgNBFInvalid
	}

	return ErrMsgUnauthorized
}

func Authenticate(ja *JWTAuth, extractor Extractor) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(rw http.ResponseWriter, req *http.Request) {
			tokenString := extractor(req)
			if tokenString == "" {
				http.Error(rw, ErrMsgNoTokenFound, http.StatusUnauthorized)
				return
			}

			ctx := req.Context()
			token, err := jwt.Parse(tokenString, ja.ParseFn(ctx, ja))
			if err != nil {
				msg := parseErrorToMsg(err)
				http.Error(rw, msg, http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(rw, ErrMsgUnauthorized, http.StatusUnauthorized)
				return
			}

			for _, validator := range ja.ValidatorFns {
				err := validator(ctx, token)
				if err != nil {
					http.Error(rw, err.Error(), http.StatusUnauthorized)
					return
				}
			}

			ctx = PassTokenToContext(ctx, token, nil)
			req = req.WithContext(ctx)

			next.ServeHTTP(rw, req)
		}

		return http.HandlerFunc(hfn)
	}
}
