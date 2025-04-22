package jwtauth

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type jwtCtxKey int

const (
	JWTCtxKey jwtCtxKey = iota
	JWTErrorCtxKey
)

func PassTokenToContext(ctx context.Context, token *jwt.Token, err error) context.Context {
	ctx = context.WithValue(ctx, JWTCtxKey, token)
	ctx = context.WithValue(ctx, JWTErrorCtxKey, err)
	return ctx
}

func GetJWTFromContext(ctx context.Context) (*jwt.Token, error) {
	err, _ := ctx.Value(JWTErrorCtxKey).(error)
	if err != nil {
		return nil, err
	}

	token, ok := ctx.Value(JWTCtxKey).(*jwt.Token)
	if !ok {
		msg := "failed to assert token"
		return nil, fmt.Errorf("%w: %s", ErrInvalidGoType, msg)
	}

	return token, nil
}
