package room

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	storage     *Storage
	redisClient *redis.Client
}

func NewService(storage *Storage, redisClient *redis.Client) *Service {
	return &Service{storage: storage, redisClient: redisClient}
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

	room, err := s.storage.CreateRoom(ctx, Room{OwnerUuid: userUuid, Name: *name, PasswordHash: passwordHash, MaxPlayers: *maxPlayers, IsPublic: *isPublic})
	if err != nil {
		return nil, err
	}
	ownerUsername, _ := s.redisClient.Get(ctx, "username:"+room.OwnerUuid.String()).Result()
	room.OwnerName = ownerUsername

	return room, nil
}
