package user

import (
	"context"
	"log"

	pb "github.com/quizverse3D/Backend/internal/pb/user"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	svc *Service
}

func NewGRPCServer(svc *Service) pb.UserServiceServer {
	return &Server{svc: svc}
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.svc.GetUser(ctx, req.GetUserId())
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return nil, err
	}

	return &pb.GetUserResponse{
		Id:       user.ID.String(),
		Username: user.Username,
	}, nil
}
