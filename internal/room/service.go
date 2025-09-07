package room

import (
	"context"
	"os"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	storage *Storage
}

func NewService(storage *Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) CreateRoom(ctx context.Context, userUuid uuid.UUID, name *string, password *string, maxPlayers *int32, isPublic *bool) (*Room, error) {
	if name == nil || *name == "" {
		return nil, ErrEmptyRoomName
	}
	if maxPlayers == nil || *maxPlayers <= 0 || *maxPlayers > 32 {
		return nil, ErrInvalidMaxPlayers
	}
	if isPublic == nil {
		return nil, ErrInvalidIsPublic
	}
	var passwordHash *string
	if password != nil {
		salt := os.Getenv("ROOMS_PASSWORD_SALT")
		combined := *password + salt
		hashByte, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashString := string(hashByte)
		passwordHash = &hashString
	}

	return nil, s.storage.CreateRoom(ctx, Room{OwnerUuid: userUuid, Name: *name, PasswordHash: passwordHash, MaxPlayers: *maxPlayers, IsPublic: *isPublic})
}
