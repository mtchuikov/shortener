package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/mtchuikov/shortener/pkg/middlewares/jwtauth"
)

const UserIDCtxKey = jwtauth.CtxKey("jwtUserIDCtxKey")

func AutoAuthenticate(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	extractor := jwtauth.ExtractFromCookie("Authorization")

	return func(next http.Handler) http.Handler {
		hfn := func(rw http.ResponseWriter, req *http.Request) {
			var valid bool
			var userID string

			ctx := req.Context()

			tokenString := extractor(req)
			if tokenString != "" {
				token, err := jwt.Parse(tokenString, ja.ParseFn(ctx, ja))
				fmt.Println(err)
				if err == nil && token.Valid {
					claims, ok := token.Claims.(jwt.MapClaims)
					if ok {
						uid, ok := claims["uid"].(string)
						if ok {
							valid = true
							userID = uid
						}
					}
				}
			}

			if !valid {
				userID = uuid.NewString()
				claims := jwt.MapClaims{"uid": userID}
				token := jwt.NewWithClaims(ja.Alg, claims)

				tokenString, err := token.SignedString(ja.SignKey)
				if err != nil {
					msg := "something went wrong"
					http.Error(rw, msg, http.StatusInternalServerError)
					return
				}

				http.SetCookie(rw,
					&http.Cookie{
						Name:  "Authorization",
						Value: tokenString,
					})
			}

			ctx = context.WithValue(ctx, UserIDCtxKey, userID)
			req = req.WithContext(ctx)

			next.ServeHTTP(rw, req)
		}

		return http.HandlerFunc(hfn)
	}
}
