package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) GetUserByID(ctx context.Context, id string) (*User, error) {
	row := s.pool.QueryRow(ctx, `SELECT id, username FROM users WHERE id = $1`, id)

	var u User
	err := row.Scan(&u.ID, &u.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &u, nil
}

func (s *Storage) CreateUser(ctx context.Context, u *User) error {
	_, err := s.pool.Exec(ctx, `INSERT INTO users (id, username) VALUES ($1, $2)`, u.ID, u.Username)
	return err
}
