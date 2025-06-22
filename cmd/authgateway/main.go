package main

import (
	"fmt"
	"log"      // стандартный логгер Go для вывода в консоль.
	"net/http" // стандартная клиент-серверная HTTP-библиотека
	"os"

	"github.com/joho/godotenv"
	"github.com/quizverse3D/Backend/internal/authgateway" // бизнес-логика
	"github.com/quizverse3D/Backend/internal/common"      // БД
)

func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system env")
	}
	// проверка наличия секрета для JWT
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	mux := http.NewServeMux() // URL-маршрутизатор

	// подключаемся к БД
	pool, err := common.NewPostgresPool(
		os.Getenv("AUTHGATEWAY_DB_USER"),
		os.Getenv("AUTHGATEWAY_DB_PASSWORD"),
		os.Getenv("AUTHGATEWAY_DB_HOST"),
		os.Getenv("AUTHGATEWAY_DB_PORT"),
		os.Getenv("AUTHGATEWAY_DB_NAME"),
	)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	redisClient, err := common.NewRedisClient()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	userService := authgateway.NewService(authgateway.NewStorage(pool), redisClient) // структура со включенным в себя Storage
	handler := authgateway.NewHandler(userService)                                   // структура-обёртка вокруг userService

	// привязка url'ов к обработчикам REST-сервиса
	mux.HandleFunc("/auth/api/v1/register", handler.Register)
	mux.HandleFunc("/auth/api/v1/login", handler.Login)
	mux.HandleFunc("/auth/api/v1/validate-token", handler.ValidateToken)
	mux.HandleFunc("/auth/api/v1/refresh-token", handler.RefreshAccessToken)

	// проверка наличия прослушиваемого REST-порта
	if os.Getenv("AUTHGATEWAY_REST_PORT") == "" {
		log.Fatal("AUTHGATEWAY_REST_PORT is not set")
	}

	restPort := fmt.Sprintf(":%s", os.Getenv("AUTHGATEWAY_REST_PORT"))

	log.Println("Authgateway REST-Service running on " + restPort)
	log.Fatal(http.ListenAndServe(restPort, mux)) // если сервер не может запуститься — log.Fatal(...) завершит программу с ошибкой и выведет сообщение
}
