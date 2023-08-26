package model

import "errors"

var (
	ErrInternalError = errors.New("internal Error")
	ErrNotFound      = errors.New("not found")
)
