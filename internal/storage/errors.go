package storage

import "errors"

const (
	ErrMsgOriginalURLNotFound      = "original url not found"
	ErrMsgOriginalURLAlreadyExists = "original url already exists"
	ErrMsgUserHasNoShortURLs       = "user has no short urls"
	ErrMsgSomethingWentWrong       = "something went wrong"
	ErrMsgNoItemsInBatch           = "no items in batch"
)

var (
	ErrOriginalURLNotFound      = errors.New(ErrMsgOriginalURLNotFound)
	ErrOriginalURLAlreadyExists = errors.New(ErrMsgOriginalURLAlreadyExists)
	ErrUserHasNoShortURLs       = errors.New(ErrMsgUserHasNoShortURLs)
	ErrNoItemsInBatch           = errors.New(ErrMsgNoItemsInBatch)
)
