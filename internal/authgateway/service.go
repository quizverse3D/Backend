package authgateway

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	storage     *Storage
	redisClient *redis.Client
	rabbitChan  *amqp.Channel
}

func NewService(storage *Storage, redisClient *redis.Client, rabbitChan *amqp.Channel) *Service {
	return &Service{storage: storage, redisClient: redisClient, rabbitChan: rabbitChan}
}

func generateSalt(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (s *Service) Register(email, password, username string) (string, error) {
	id := uuid.NewString()

	salt, err := generateSalt(16)
	if err != nil {
		return "", err
	}
	combined := password + salt
	hashed, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	u := Auth{
		ID:           id,
		Email:        email,
		Password:     string(hashed),
		PasswordSalt: salt,
	}

	err = s.storage.CreateAuth(u)
	if err != nil {
		return "", err
	}

	err = s.rabbitChan.Publish("", "user_registered", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(fmt.Sprintf(`{"userId":"%s","userName":"%s"}`, id, username)),
	})
	if err != nil {
		return "", fmt.Errorf("failed to publish event: %w", err)
	}

	return id, nil
}

func (s *Service) Login(email, password string) (string, string, error) {
	u, ok := s.storage.GetAuth(email)
	if !ok {
		return "", "", ErrInvalidCreds
	}

	combined := password + u.PasswordSalt

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

func (s *Service) RefreshAccessToken(refreshToken string) (string, error) {
	userID, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	key := fmt.Sprintf("refresh:%s", userID)

	stored, err := s.redisClient.Get(ctx, key).Result()
	if err != nil || stored != refreshToken {
		return "", ErrInvalidCreds
	}

	return GenerateAccessToken(userID)
}

func (s *Service) UpdatePassword(uuid, newPassword, oldPassword string) error {
	u, err := s.storage.GetCredInfoByUuid(uuid)
	if err != nil {
		return err
	}

	salt := u.PasswordSalt
	combined := oldPassword + salt

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(combined))
	if err != nil {
		return ErrInvalidPassword
	}

	salt, err = generateSalt(16)
	if err != nil {
		return err
	}
	combined = newPassword + salt
	newHashed, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.storage.UpdatePasswordForUuid(uuid, string(newHashed), salt)
	if err != nil {
		return err
	}

	return nil
}
