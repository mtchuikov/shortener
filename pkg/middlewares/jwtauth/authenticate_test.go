package jwtauth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
)

type testAuthenticateSuite struct {
	suite.Suite
	secretKey []byte
	ja        *JWTAuth
	handler   http.Handler
	rr        *httptest.ResponseRecorder
	req       *http.Request
}

func TestAuthenticate(t *testing.T) {
	suite.Run(t, new(testAuthenticateSuite))
}

func (s *testAuthenticateSuite) SetupTest() {
	s.secretKey = []byte("jwt")
	s.ja = New(s.secretKey)

	middleware := Authenticate(s.ja, ExtractFromAuthHeader)
	s.handler = middleware(http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusCreated)
		},
	))

	s.rr = httptest.NewRecorder()
	s.req = httptest.NewRequest(http.MethodGet, "/", nil)
}

func (s *testAuthenticateSuite) TestAuthenticate_Success() {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, _ := token.SignedString(s.secretKey)

	s.req.Header.Set("Authorization", "Bearer "+tokenString)

	s.handler.ServeHTTP(s.rr, s.req)
	s.Require().Equal(http.StatusCreated, s.rr.Code)
}

func (s *testAuthenticateSuite) TestAuthenticate_NoToken() {
	s.handler.ServeHTTP(s.rr, s.req)
	s.Require().Equal(http.StatusUnauthorized, s.rr.Code)

	payload, _ := io.ReadAll(s.rr.Body)
	msg := string(payload)
	msg = strings.ReplaceAll(msg, "\n", "")

	s.Require().Equal(ErrMsgNoTokenFound, msg)
}

func (s *testAuthenticateSuite) TestAuthenticate_ExpiredToken() {
	claims := jwt.MapClaims{"exp": jwt.NewNumericDate(time.Now().Add(-time.Hour))}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(s.secretKey)
	s.req.Header.Set("Authorization", "Bearer "+tokenString)

	s.handler.ServeHTTP(s.rr, s.req)
	s.Require().Equal(http.StatusUnauthorized, s.rr.Code)

	payload, _ := io.ReadAll(s.rr.Body)
	msg := string(payload)
	msg = strings.ReplaceAll(msg, "\n", "")

	s.Require().Equal(ErrMsgExpired, msg)

}
