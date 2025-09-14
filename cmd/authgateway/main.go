package main

import (
	"fmt"
	"log"      // стандартный логгер Go для вывода в консоль.
	"net/http" // стандартная клиент-серверная HTTP-библиотека
	"os"

	"github.com/joho/godotenv"
	"github.com/quizverse3D/Backend/internal/authgateway" // бизнес-логика
	"github.com/quizverse3D/Backend/internal/common"      // БД
	"github.com/streadway/amqp"
)

func main() {
	// .env
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system env")
	}

	mux := http.NewServeMux()

	// PostgreSQL
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

	// Redis
	redisClient, err := common.NewRedisClient()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// RabbitMQ
	rabbitConn, err := amqp.Dial(os.Getenv(("RABBITMQ_URL")))
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	rabbitChan, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("failed to open RabbitMQ channel: %v", err)
	}
	defer rabbitChan.Close()

	// Service and Storage
	authService := authgateway.NewService(authgateway.NewStorage(pool), redisClient, rabbitChan) // структура со включенным в себя Storage
	handler := authgateway.NewHandler(authService)                                               // структура-обёртка вокруг authService

	// привязка url'ов к обработчикам REST-сервиса
	mux.HandleFunc("/auth/api/v1/register", handler.Register)
	mux.HandleFunc("/auth/api/v1/login", handler.Login)
	mux.HandleFunc("/auth/api/v1/validate-token", handler.ValidateToken)
	mux.HandleFunc("/auth/api/v1/refresh-token", handler.RefreshAccessToken)
	mux.HandleFunc("/auth/api/v1/update-password", handler.UpdatePassword)

	// привязка gRPC-сервисов для маршрутизации
	grpcUserAddr := fmt.Sprintf("%s:%s", os.Getenv("USERS_GRPC_HOST"), os.Getenv("USERS_GRPC_PORT"))
	userRestPrefix := "/user/api/v1/"
	userRoute, err := authgateway.NewUserGrpcServiceRoute(grpcUserAddr, userRestPrefix)
	if err != nil {
		log.Fatalf("failed to create userRoute: %v", err)
	}
	defer userRoute.Conn.Close()
	mux.Handle(userRestPrefix, authgateway.AuthMiddleWare(authgateway.ProxyHandler(userRoute)))

	grpcRoomAddr := fmt.Sprintf("%s:%s", os.Getenv("ROOMS_GRPC_HOST"), os.Getenv("ROOMS_GRPC_PORT"))
	roomRestPrefix := "/room/api/v1/"
	roomRoute, err := authgateway.NewRoomGrpcServiceRoute(grpcRoomAddr, roomRestPrefix)
	if err != nil {
		log.Fatalf("failed to create roomRoute: %v", err)
	}
	defer roomRoute.Conn.Close()
	mux.Handle(roomRestPrefix, authgateway.AuthMiddleWare(authgateway.ProxyHandler(roomRoute)))

	// REST Server listening (в конце)
	restPort := fmt.Sprintf(":%s", os.Getenv("AUTHGATEWAY_REST_PORT"))
	log.Println("Authgateway REST-Service running on " + restPort)
	log.Fatal(http.ListenAndServe(restPort, mux))
}
