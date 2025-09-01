package authgateway

import "errors"

var (
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidCreds    = errors.New("invalid credentials")
	ErrInvalidPassword = errors.New("invalid password")
)
