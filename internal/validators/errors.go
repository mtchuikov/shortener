package validators

import "errors"

const (
	ErrMsgInvalidCorrelationID = "invalid correlation id"
	ErrMsgInvalidOriginalURL   = "invalid original url"
	ErrMsgInvalidSlug          = "invalid slug"
	ErrMsgInvalidUserID        = "invalid user id"
)

var (
	ErrInvalidCorrelationID = errors.New(ErrMsgInvalidCorrelationID)
	ErrInvalidOriginalURL   = errors.New(ErrMsgInvalidOriginalURL)
	ErrInvalidSlug          = errors.New(ErrMsgInvalidSlug)
	ErrInvalidUserID        = errors.New(ErrMsgInvalidUserID)
)
