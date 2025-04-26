package jwtauth

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type CtxKey string

const (
	TokenCtxKey = CtxKey("jwtTokenCtxKey")
	ErrorCtxKey = CtxKey("jwtErrorCtxKey")
)

func PassTokenToContext(ctx context.Context, token *jwt.Token, err error) context.Context {
	ctx = context.WithValue(ctx, TokenCtxKey, token)
	ctx = context.WithValue(ctx, ErrorCtxKey, err)
	return ctx
}

func GetTokenFromContext(ctx context.Context) (*jwt.Token, error) {
	err, _ := ctx.Value(ErrorCtxKey).(error)
	if err != nil {
		return nil, err
	}

	token, ok := ctx.Value(TokenCtxKey).(*jwt.Token)
	if !ok {
		msg := "failed to assert token"
		return nil, fmt.Errorf("%w: %s", ErrInvalidGoType, msg)
	}

	return token, nil
}
