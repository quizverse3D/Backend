package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserParamsNotFound = errors.New("user client params not found")
)
