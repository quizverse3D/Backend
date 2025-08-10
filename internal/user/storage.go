package user

import (
	"context"

	"github.com/google/uuid"
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
	_, usersErr := s.pool.Exec(ctx, `INSERT INTO users (id, username) VALUES ($1, $2)`, u.ID, u.Username)
	if usersErr != nil {
		return usersErr
	}
	_, paramsErr := s.pool.Exec(ctx, `INSERT INTO params (user_uuid) VALUES ($1)`, u.ID)
	return paramsErr
}

func (s *Storage) GetUserClientParamsByUuid(ctx context.Context, uuid uuid.UUID) (*ClientParams, error) {
	row := s.pool.QueryRow(ctx, `SELECT user_uuid, lang_code, sound_volume, game_sound_enabled FROM params WHERE user_uuid = $1`, uuid)

	var p ClientParams
	err := row.Scan(&p.UserUuid, &p.LangCode, &p.SoundVolume, &p.IsGameSoundEnabled)
	if err != nil {
		return nil, ErrUserParamsNotFound
	}

	return &p, nil
}
