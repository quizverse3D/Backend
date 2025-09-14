package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	storage     *Storage
	redisClient *redis.Client
}

func NewService(storage *Storage, redisClient *redis.Client) *Service {
	return &Service{storage: storage, redisClient: redisClient}
}

func (s *Service) GetUser(ctx context.Context, userID string) (*User, error) {
	return s.storage.GetUserByID(ctx, userID)
}

func (s *Service) CreateUser(ctx context.Context, u *User) error {
	return s.storage.CreateUser(ctx, u)
}

func (s *Service) GetUserClientParamsByUuid(ctx context.Context, userUuid uuid.UUID) (*ClientParams, error) {
	return s.storage.GetUserClientParamsByUuid(ctx, userUuid)
}

func (s *Service) SetUserClientParamsByUuid(ctx context.Context, userUuid uuid.UUID,
	langCode *string,
	soundVolume *int32,
	isGameSoundEnabled *bool) (*ClientParams, error) {
	// validate
	if langCode != nil {
		validLangCodes := map[string]struct{}{
			"EN": {},
			"RU": {},
		}
		if _, ok := validLangCodes[*langCode]; !ok {
			return nil, ErrUserParamsInvalidLangCode
		}
	}
	if soundVolume != nil && (*soundVolume < 0 || *soundVolume > 100) {
		return nil, ErrUserParamsInvalidSoundVolume
	}
	return s.storage.SetUserClientParamsByUuid(ctx, userUuid, langCode, soundVolume, isGameSoundEnabled)
}

func (s *Service) SyncUsernamesToRedis(ctx context.Context, userUuid *uuid.UUID) error {
	if userUuid != nil {
		// один пользователь
		user, err := s.storage.GetUserByID(ctx, userUuid.String())
		if err != nil {
			return err
		}
		return s.redisClient.Set(ctx, "username:"+user.ID.String(), user.Username, 0).Err()
	}

	// все пользователи
	users, err := s.storage.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		if err := s.redisClient.Set(ctx, "username:"+user.ID.String(), user.Username, 0).Err(); err != nil {
			return err
		}
	}
	return nil
}
