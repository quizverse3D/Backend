package user

import (
	"context"
	"fmt"
	"os"
	"time"

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

func (s *Service) Register(username, password string) (string, error) {
	id := uuid.NewString()

	salt := os.Getenv("PASSWORD_SALT")
	combined := password + salt
	hashed, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	u := User{
		ID:            id,
		Username:      username,
		Password:      string(hashed),
		HashAlgorithm: "bcrypt",
	}

	err = s.storage.CreateUser(u)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) Login(username, password string) (string, string, error) {
	u, ok := s.storage.GetUser(username)
	if !ok {
		return "", "", ErrInvalidCreds
	}

	salt := os.Getenv("PASSWORD_SALT")
	combined := password + salt

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(combined))
	if err != nil {
		return "", "", ErrInvalidCreds
	}

	accessToken, err := GenerateAccessToken(u.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(u.ID)
	if err != nil {
		return "", "", err
	}

	ctx := context.Background()
	key := fmt.Sprintf("refresh:%s", u.ID)
	err = s.redisClient.Set(ctx, key, refreshToken, 7*24*time.Hour).Err()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) ValidateAccessToken(tokenStr string) (string, error) {
	return ValidateAccessToken(tokenStr)
}
