package room

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

func (s *Storage) CreateRoom(ctx context.Context, room Room) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO rooms (owner_id, name, password_hash, max_players, is_public)
		VALUES ($1, $2, $3, $4, $5)`, room.OwnerUuid, room.Name, room.PasswordHash, room.MaxPlayers, room.IsPublic)
	return err
}
