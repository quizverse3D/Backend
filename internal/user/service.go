package user

import (
	"context"
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
