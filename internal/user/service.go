package user

import (
	"github.com/google/uuid"
)

type Service struct {
	storage *Storage
}

func NewService() *Service {
	return &Service{storage: NewStorage()}
}

func (s *Service) Register(username, password string) (string, error) {
	id := uuid.NewString()
	u := User{
		ID:       id,
		Username: username,
		Password: password,
	}

	err := s.storage.CreateUser(u)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) Login(username, password string) (string, error) {
	u, ok := s.storage.GetUser(username)
	if !ok || u.Password != password {
		return "", ErrInvalidCreds
	}

	return GenerateJWT(u.ID)
}
