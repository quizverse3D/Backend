package user

import "errors"

var (
    ErrUserExists    = errors.New("user already exists")
    ErrInvalidCreds  = errors.New("invalid credentials")
)
