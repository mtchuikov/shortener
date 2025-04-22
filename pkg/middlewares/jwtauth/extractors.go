package jwtauth

import (
	"net/http"
	"strings"
)

type TokenExtractor func(req *http.Request) string

func ExtractFromAuthHeader(req *http.Request) string {
	bearer := req.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToLower(bearer[0:7]) == "bearer " {
		return bearer[7:]
	}

	return ""
}

func ExtractFromCookie(cookie string) TokenExtractor {
	return func(req *http.Request) string {
		cookie, err := req.Cookie(cookie)
		if err == http.ErrNoCookie {
			return ""
		}

		return cookie.Value
	}
}

func ExtractFromQuery(param string) TokenExtractor {
	return func(req *http.Request) string {
		return req.URL.Query().Get(param)
	}
}

func MultiTokenExtractor(extractors ...TokenExtractor) TokenExtractor {
	return func(req *http.Request) string {
		for _, extractor := range extractors {
			tokenString := extractor(req)
			if tokenString != "" {
				return tokenString
			}
		}

		return ""
	}
}
