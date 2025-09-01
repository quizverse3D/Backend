package authgateway

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

func (s *Storage) CreateAuth(u Auth) error {
	_, err := s.db.Exec(context.Background(),
		"INSERT INTO credentials (id, email, password, hash_algorithm) VALUES ($1, $2, $3, $4)",
		u.ID, u.Email, u.Password, u.HashAlgorithm,
	)

	if err == nil {
		// ошибок нет, пользователь создан
		return nil
	}

	return err
}

func (s *Storage) GetAuth(email string) (Auth, bool) {
	row := s.db.QueryRow(context.Background(),
		"SELECT id, email, password, hash_algorithm FROM credentials WHERE email = $1",
		email,
	)

	var u Auth
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.HashAlgorithm)
	if err != nil {
		return Auth{}, false
	}

	return u, true
}

func (s *Storage) GetCredInfoByUuid(uuid string) (Auth, error) {
	row := s.db.QueryRow(context.Background(), "SELECT id, email, password, hash_algorithm FROM credentials WHERE id = $1", uuid)

	var u Auth
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.HashAlgorithm)
	if err != nil {
		return Auth{}, err
	}

	return u, nil
}

func (s *Storage) UpdatePasswordForUuid(uuid, password, hashAlgorithm string) error {
	_, err := s.db.Exec(context.Background(), "UPDATE credentials SET password = $1, hash_algorithm = $2 WHERE id = $3", password, hashAlgorithm, uuid)

	if err == nil {
		return nil
	}

	return err
}
