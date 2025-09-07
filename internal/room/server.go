package room

import (
	pb "github.com/quizverse3D/Backend/internal/pb/room"
)

type Server struct {
	pb.UnimplementedRoomServiceServer
	svc *Service
}

func NewGRPCServer(svc *Service) pb.RoomServiceServer {
	return &Server{svc: svc}
}
