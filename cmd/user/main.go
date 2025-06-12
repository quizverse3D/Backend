package main

import (
	"log"      // стандартный логгер Go для вывода в консоль.
	"net/http" // стандартная клиент-серверная HTTP-библиотека
	"os"

	"github.com/joho/godotenv"
	"github.com/quizverse3D/Backend/internal/user" // бизнес-логика
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

	userService := user.NewService()        // структура со включенным в себя Storage
	handler := user.NewHandler(userService) // структура-обёртка вокруг userService

	// привязка url'ов к обработчикам сервиса
	mux.HandleFunc("/api/v1/register", handler.Register)
	mux.HandleFunc("/api/v1/login", handler.Login)

	log.Println("User Service running on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux)) // если сервер не может запуститься — log.Fatal(...) завершит программу с ошибкой и выведет сообщение
}
