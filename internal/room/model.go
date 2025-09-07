package room

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID           uuid.UUID
	OwnerUuid    uuid.UUID
	Name         string
	PasswordHash *string
	MaxPlayers   int32
	CreatedAt    *time.Time
	IsPublic     bool
}
