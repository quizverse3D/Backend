package common

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool() (*pgxpool.Pool, error) {
	// собираем строку подключения
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("USERS_DB_USER"),
		os.Getenv("USERS_DB_PASSWORD"),
		os.Getenv("USERS_DB_HOST"),
		os.Getenv("USERS_DB_PORT"),
		os.Getenv("USERS_DB_NAME"),
	)

	// Создаём контекст с таймаутом 5 секунд (если БД не отвечает — не будем ждать вечно)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // обязательно вызвать, чтобы освободить ресурсы

	// Инициализируем пул соединений
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	// пингуем базу (проверка, что соединение действительно рабочее)
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
