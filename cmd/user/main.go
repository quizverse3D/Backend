package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/quizverse3D/Backend/internal/common"
	pb "github.com/quizverse3D/Backend/internal/pb/user"
	"github.com/quizverse3D/Backend/internal/user"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

func main() {
	// .env
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system env")
	}

	// PostgreSQL
	pool, err := common.NewPostgresPool(
		os.Getenv("USERS_DB_USER"),
		os.Getenv("USERS_DB_PASSWORD"),
		os.Getenv("USERS_DB_HOST"),
		os.Getenv("USERS_DB_PORT"),
		os.Getenv("USERS_DB_NAME"),
	)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
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
	storage := user.NewStorage(pool)
	service := user.NewService(storage)

	// gRPC Server
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, user.NewGRPCServer(service))

	listener, err := net.Listen("tcp", ":"+os.Getenv("USERS_GRPC_PORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// регистрация rabbitmq consumer'ов
	consumers := []common.Consumer{
		*common.NewConsumer(rabbitChan, "user_registered", user.UserRegisteredHandler(service)),
	}

	for _, c := range consumers {
		if err := c.DeclareQueue(); err != nil {
			log.Fatalf("failed to declare queue: %v", err)
		}
		if err := c.Listen(context.Background()); err != nil {
			log.Fatalf("failed to start consumer: %v", err)
		}
	}

	// Выполняется в конце, прослушивание gRPC
	log.Println("User gRPC-service running on " + os.Getenv("USERS_GRPC_PORT"))
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
