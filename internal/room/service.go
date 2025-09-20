package room

import (
	"context"
	"crypto/rand"
	"encoding/base64"

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

func generateSalt(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
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
	var passwordSalt string
	if password != nil {
		ps, err := generateSalt(12)
		if err != nil {
			return nil, err
		}
		passwordSalt = ps
		combined := *password + passwordSalt
		hashByte, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashString := string(hashByte)
		passwordHash = &hashString
	}

	room, err := s.storage.CreateRoom(ctx, Room{OwnerUuid: userUuid, Name: *name, PasswordHash: passwordHash, PasswordSalt: passwordSalt, MaxPlayers: *maxPlayers, IsPublic: *isPublic})
	if err != nil {
		return nil, err
	}
	ownerUsername, _ := s.redisClient.Get(ctx, "username:"+room.OwnerUuid.String()).Result()
	room.OwnerName = ownerUsername

	return room, nil
}
