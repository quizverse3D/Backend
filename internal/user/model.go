package user

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Username string
}

type ClientParams struct {
	UserUuid           uuid.UUID
	LangCode           string
	SoundVolume        int16
	IsGameSoundEnabled bool
}
