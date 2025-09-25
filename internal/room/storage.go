package room

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

func (s *Storage) CreateRoom(ctx context.Context, room Room) (*Room, error) {
	row := s.pool.QueryRow(ctx, `
		INSERT INTO rooms (owner_id, name, password_hash, password_salt, max_players, is_public)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, owner_id, name, max_players, created_at, is_public
		`, room.OwnerUuid, room.Name, room.PasswordHash, room.PasswordSalt, room.MaxPlayers, room.IsPublic)

	var r Room
	if err := row.Scan(&r.ID, &r.OwnerUuid, &r.Name, &r.MaxPlayers, &r.CreatedAt, &r.IsPublic); err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *Storage) GetRoomById(ctx context.Context, uuid uuid.UUID) (*Room, error) {
	row := s.pool.QueryRow(ctx, `SELECT id, name, owner_id, password_hash, password_salt, max_players, created_at, is_public FROM rooms WHERE id = $1`, uuid.String())
	var r Room
	if err := row.Scan(&r.ID, &r.Name, &r.OwnerUuid, &r.PasswordHash, &r.PasswordSalt, &r.MaxPlayers, &r.CreatedAt, &r.IsPublic); err != nil {
		return nil, ErrRoomNotFound
	}
	return &r, nil
}
