package user

import (
	"context"
	"fmt"
	"strings"

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

func (s *Storage) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, username FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
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

func (s *Storage) SetUserClientParamsByUuid(ctx context.Context, uuid uuid.UUID, langCode *string,
	soundVolume *int32,
	isGameSoundEnabled *bool) (*ClientParams, error) {

	setParts := []string{}
	args := []interface{}{}
	argPos := 1

	if langCode != nil {
		setParts = append(setParts, fmt.Sprintf("lang_code = $%d", argPos))
		args = append(args, *langCode)
		argPos++
	}
	if soundVolume != nil {
		setParts = append(setParts, fmt.Sprintf("sound_volume = $%d", argPos))
		args = append(args, *soundVolume)
		argPos++
	}
	if isGameSoundEnabled != nil {
		setParts = append(setParts, fmt.Sprintf("game_sound_enabled = $%d", argPos))
		args = append(args, *isGameSoundEnabled)
		argPos++
	}

	if len(setParts) == 0 {
		return s.GetUserClientParamsByUuid(ctx, uuid)
	}

	query := fmt.Sprintf("UPDATE params SET %s WHERE user_uuid = $%d RETURNING user_uuid, lang_code, sound_volume, game_sound_enabled", strings.Join(setParts, ", "), argPos)
	args = append(args, uuid)
	row := s.pool.QueryRow(ctx, query, args...)

	var clientParams ClientParams
	err := row.Scan(&clientParams.UserUuid, &clientParams.LangCode, &clientParams.SoundVolume, &clientParams.IsGameSoundEnabled)
	if err != nil {
		return nil, err
	}

	return &clientParams, nil
}
