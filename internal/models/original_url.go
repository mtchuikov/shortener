package models

import (
	"errors"
	"net/url"
	"strings"
)

type OriginalURL string

var ErrInvalidOriginalURL = errors.New("invalid original url")

func validateOriginalURL(u string) error {
	parsed, err := url.Parse(u)
	if err != nil {
		return ErrInvalidOriginalURL
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrInvalidOriginalURL
	}

	if parsed.Host == "" || strings.IndexRune(parsed.Host, '\\') > 0 {
		return ErrInvalidOriginalURL
	}

	return nil
}

func NewOriginalURL(u string) (OriginalURL, error) {
	err := validateOriginalURL(u)
	if err != nil {
		return "", err
	}

	return OriginalURL(u), nil
}

func (o OriginalURL) String() string {
	return string(o)
}
