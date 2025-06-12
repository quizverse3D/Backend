package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{db: pool}
}

func (s *Storage) CreateUser(u User) error {
	_, err := s.db.Exec(context.Background(),
		"INSERT INTO users (id, username, password) VALUES ($1, $2, $3)",
		u.ID, u.Username, u.Password,
	)

	if err == nil {
		// ошибок нет, пользователь создан
		return nil
	}

	return err
}

func (s *Storage) GetUser(username string) (User, bool) {
	row := s.db.QueryRow(context.Background(),
		"SELECT id, username, password FROM users WHERE username = $1",
		username,
	)

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return User{}, false
	}

	return u, true
}
