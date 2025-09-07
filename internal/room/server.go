package room

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	pb "github.com/quizverse3D/Backend/internal/pb/room"
)

type Server struct {
	pb.UnimplementedRoomServiceServer
	svc *Service
}

func NewGRPCServer(svc *Service) pb.RoomServiceServer {
	return &Server{svc: svc}
}

func (s *Server) CreateRoom(ctx context.Context, req *pb.CreateRoomParamsRequest) (*pb.CreateRoomParamsResponse, error) {
	// parse pb
	userUuid, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return nil, fmt.Errorf("invalid user_uuid: %w", err)
	}

	// call service
	// TODO: после реализации метода получения комнаты по uuid, возвращать результат создания
	_, err = s.svc.CreateRoom(ctx, userUuid, &req.Name, req.Password, &req.MaxPlayers, &req.IsPublic)
	if err != nil {
		log.Printf("failed to create room: %v", err)
		return nil, err
	}

	return &pb.CreateRoomParamsResponse{}, nil
}
