package jwtauth

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
)

type testContextSuite struct {
	suite.Suite
	ctx   context.Context
	token *jwt.Token
}

func TestContext(t *testing.T) {
	suite.Run(t, new(testContextSuite))
}

func (s *testContextSuite) SetupTest() {
	s.ctx = context.Background()
	s.token = &jwt.Token{}
}

func (s *testContextSuite) TestPassAndGetToken_Success() {
	claims := jwt.MapClaims{"foo": "bar", "nbf": float64(123)}
	s.token.Claims = claims

	ctx := PassTokenToContext(s.ctx, s.token, nil)
	gotToken, err := GetTokenFromContext(ctx)

	s.Require().NoError(err)
	s.Require().Equal(s.token, gotToken)
}

func (s *testContextSuite) TestGetTokenFromContext_NoToken() {
	gotToken, err := GetTokenFromContext(s.ctx)

	s.Require().Nil(gotToken)
	s.Require().ErrorIs(err, ErrInvalidGoType)
}

func (s *testContextSuite) TestGetTokenFromContext_InvalidTokenType() {
	ctx := context.WithValue(s.ctx, TokenCtxKey, "not a token")
	gotToken, err := GetTokenFromContext(ctx)

	s.Require().Nil(gotToken)
	s.Require().ErrorIs(err, ErrInvalidGoType)
}

func (s *testContextSuite) TestGetTokenFromContext_PreservesError() {
	sampleErr := errors.New("sample failure")
	s.ctx = PassTokenToContext(s.ctx, nil, sampleErr)

	gotToken, err := GetTokenFromContext(s.ctx)

	s.Require().Nil(gotToken)
	s.Require().ErrorIs(err, sampleErr, err)
}
