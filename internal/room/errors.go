package room

import "errors"

var (
	ErrEmptyRoomName     = errors.New("Empty room name is not allowed")
	ErrInvalidMaxPlayers = errors.New("Max players value must be 1-32")
	ErrInvalidIsPublic   = errors.New("Public visibility parameter must be True or False")
)
