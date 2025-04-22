package models

import (
	"errors"
	"regexp"
)

type ShortenID string

var (
	validShortenID      = regexp.MustCompile(`^[A-Za-z0-9]{8,12}$`)
	ErrInvalidShortenID = errors.New("invalid shorten id")
)

func NewShortenID(id string) (ShortenID, error) {
	valid := validShortenID.MatchString(id)
	if !valid {
		return "", ErrInvalidShortenID
	}

	return ShortenID(id), nil
}

func (s ShortenID) String() string {
	return string(s)
}
