package user

import (
	"context"
	"log"

	"github.com/google/uuid"
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

func (s *Server) GetUserClientParams(ctx context.Context, req *pb.GetUserClientParamsRequest) (*pb.GetUserClientParamsResponse, error) {
	params, err := s.svc.GetUserClientParamsByUuid(ctx, uuid.MustParse(req.GetUserUuid()))
	if err != nil {
		log.Printf("failed to get user client params: %v", err)
		return nil, err
	}

	return &pb.GetUserClientParamsResponse{
		UserUuid:           params.UserUuid.String(),
		LangCode:           params.LangCode,
		SoundVolume:        int32(params.SoundVolume),
		IsGameSoundEnabled: params.IsGameSoundEnabled,
	}, nil
}
