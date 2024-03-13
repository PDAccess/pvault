package domain

import "errors"

var (
	ErrNotFound     = errors.New("no data found")
	ErrInvalidToken = errors.New("invalid token")
)
