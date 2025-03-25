package handlers

import "errors"

var (
	errFailedToReadBody      = errors.New("failed to read body")
	errFailedToMarshalJSON   = errors.New("failed to marshal json")
	errFailedToUnmarshalJSON = errors.New("failed to unmarshal json")
)
