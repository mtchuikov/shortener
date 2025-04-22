package jwtauth

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
)

type testUtilsSuite struct {
	suite.Suite
	ctx   context.Context
	token *jwt.Token
}

func TestUtils(t *testing.T) {
	suite.Run(t, new(testUtilsSuite))
}

func (s *testUtilsSuite) SetupTest() {
	s.ctx = context.Background()
	s.token = &jwt.Token{}
}

func (s *testUtilsSuite) TestPassAndGetToken_Success() {
	claims := jwt.MapClaims{"foo": "bar", "nbf": float64(123)}
	s.token.Claims = claims

	ctx := PassTokenToContext(s.ctx, s.token, nil)
	gotToken, err := GetJWTFromContext(ctx)

	s.Require().NoError(err)
	s.Require().Equal(s.token, gotToken)
}

func (s *testUtilsSuite) TestGetJWTFromContext_NoToken() {
	gotToken, err := GetJWTFromContext(s.ctx)

	s.Require().Nil(gotToken)
	s.Require().ErrorIs(err, ErrInvalidGoType)
}

func (s *testUtilsSuite) TestGetJWTFromContext_InvalidTokenType() {
	ctx := context.WithValue(s.ctx, JWTCtxKey, "not a token")
	gotToken, err := GetJWTFromContext(ctx)

	s.Require().Nil(gotToken)
	s.Require().ErrorIs(err, ErrInvalidGoType)
}

func (s *testUtilsSuite) TestPassTokenToContext_PreservesError() {
	sampleErr := errors.New("sample failure")
	s.ctx = PassTokenToContext(s.ctx, nil, sampleErr)

	gotToken, err := GetJWTFromContext(s.ctx)

	s.Require().Nil(gotToken)
	s.Require().ErrorIs(err, sampleErr, err)
}
