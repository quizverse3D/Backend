package room

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID           uuid.UUID
	OwnerUuid    uuid.UUID
	OwnerName    string
	Name         string
	PasswordHash *string
	PasswordSalt string
	MaxPlayers   int32
	CreatedAt    *time.Time
	IsPublic     bool
}
