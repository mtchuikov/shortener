package jwtauth

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type testExtractorsSuite struct {
	suite.Suite
	tokenString string
}

func TestExtractors(t *testing.T) {
	suite.Run(t, new(testExtractorsSuite))
}

func (s *testExtractorsSuite) SetupTest() {
	s.tokenString = "jwt"
}

func (s *testExtractorsSuite) TestExtractFromAuthHeader_Success() {
	req := &http.Request{Header: make(http.Header)}
	req.Header.Add("Authorization", "Bearer jwt")

	tokenString := ExtractFromAuthHeader(req)
	s.Equal(s.tokenString, tokenString)
}

func (s *testExtractorsSuite) TestExtractFromAuthHeader_InvalidHeader() {
	s.tokenString = ""

	req := &http.Request{Header: make(http.Header)}
	req.Header.Add("Authorization", "jwt")

	tokenString := ExtractFromAuthHeader(req)
	s.Equal(s.tokenString, tokenString)
}

func (s *testExtractorsSuite) TestExtractFromQuery_Success() {
	u := &url.URL{}
	q := u.Query()
	q.Set("token", s.tokenString)
	u.RawQuery = q.Encode()

	req := &http.Request{URL: u}
	extractor := ExtractFromQuery("token")

	tokenString := extractor(req)
	s.Equal(s.tokenString, tokenString)
}

func (s *testExtractorsSuite) TestExtractFromQuery_NotFound() {
	u := &url.URL{}
	req := &http.Request{URL: u}
	extractor := ExtractFromQuery("token")

	tokenString := extractor(req)
	s.Empty(tokenString)
}

func (s *testExtractorsSuite) TestExtractFromCookie_Success() {
	req := &http.Request{Header: make(http.Header)}
	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: s.tokenString,
	})
	extractor := ExtractFromCookie("jwt")

	tokenString := extractor(req)
	s.Equal(s.tokenString, tokenString)
}

func (s *testExtractorsSuite) TestExtractFromCookie_NoCookie() {
	req := &http.Request{Header: make(http.Header)}
	extractor := ExtractFromCookie("jwt")

	tokenString := extractor(req)
	s.Empty(tokenString)
}
