package user

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	storage *Storage
}

func NewService(storage *Storage) *Service {
	return &Service{storage: storage}
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
