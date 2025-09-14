package room

import "errors"

var (
	ErrEmptyRoomName     = errors.New("empty room name is not allowed")
	ErrInvalidMaxPlayers = errors.New("max players value must be 1-32")
	ErrInvalidIsPublic   = errors.New("public visibility parameter must be True or False")
)
