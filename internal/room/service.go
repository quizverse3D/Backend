package room

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"

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

func (s *Service) GetRoomById(ctx context.Context, uuid uuid.UUID, isHiddenPassword bool) (*Room, error) {
	room, err := s.storage.GetRoomById(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if isHiddenPassword {
		room.PasswordHash = nil
		room.PasswordSalt = ""
	}
	room.OwnerName, err = s.redisClient.Get(ctx, "username:"+room.OwnerUuid.String()).Result()
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (s *Service) SearchRooms(ctx context.Context, search *string, page, size int32) ([]Room, int64, error) {
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	if page <= 0 {
		page = 1
	}

	var normalizedValue string
	var normalized *string
	if search != nil {
		trimmed := strings.TrimSpace(*search)
		if trimmed != "" {
			normalizedValue = trimmed
			normalized = &normalizedValue
		}
	}

	offset := (page - 1) * size
	rooms, total, err := s.storage.SearchRooms(ctx, normalized, size, offset)
	if err != nil {
		return nil, 0, err
	}

	for i := range rooms {
		username, err := s.redisClient.Get(ctx, "username:"+rooms[i].OwnerUuid.String()).Result()
		if err != nil {
			username = ""
		}
		if username != "" {
			rooms[i].OwnerName = username
		}
		rooms[i].PasswordHash = nil
		rooms[i].PasswordSalt = ""
	}

	return rooms, total, nil
}

func (s *Service) DeleteRoom(ctx context.Context, userUuid, roomUuid uuid.UUID) error {
	room, err := s.storage.GetRoomById(ctx, roomUuid)
	if err != nil {
		return err
	}
	if room.OwnerUuid != userUuid {
		return ErrRoomForbidden
	}
	return s.storage.DeleteRoom(ctx, roomUuid)
}
