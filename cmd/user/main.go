package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/quizverse3D/Backend/internal/common"
	pb "github.com/quizverse3D/Backend/internal/pb/user"
	"github.com/quizverse3D/Backend/internal/user"
	"google.golang.org/grpc"
)

func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system env")
	}

	grpcPort := os.Getenv("USERS_GRPC_PORT")
	if grpcPort == "" {
		log.Fatal("USERS_GRPC_PORT is not set")
	}

	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

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

	storage := user.NewStorage(pool)
	service := user.NewService(storage)

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, user.NewGRPCServer(service))

	log.Println("User gRPC-service running on " + grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
