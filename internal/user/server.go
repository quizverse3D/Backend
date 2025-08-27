package user

import (
	"context"
	"fmt"
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
		SoundVolume:        &params.SoundVolume,
		IsGameSoundEnabled: &params.IsGameSoundEnabled,
	}, nil
}

func (s *Server) SetUserClientParams(ctx context.Context, req *pb.SetUserClientParamsRequest) (*pb.SetUserClientParamsResponse, error) {
	// parse pb
	userUuid, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return nil, fmt.Errorf("invalid user_uuid: %w", err)
	}
	// * to disable goland type auto default values
	var langCode *string
	if req.LangCode != nil {
		langCode = &req.LangCode.Value
	}

	var soundVolume *int32
	if req.SoundVolume != nil {
		soundVolume = &req.SoundVolume.Value
	}

	var isGameSoundEnabled *bool
	if req.IsGameSoundEnabled != nil {
		isGameSoundEnabled = &req.IsGameSoundEnabled.Value
	}

	// call service
	params, err := s.svc.SetUserClientParamsByUuid(ctx, userUuid, langCode, soundVolume, isGameSoundEnabled)
	if err != nil {
		log.Printf("failed to set user client params: %v", err)
		return nil, err
	}

	return &pb.SetUserClientParamsResponse{
		UserUuid:           params.UserUuid.String(),
		LangCode:           params.LangCode,
		SoundVolume:        &params.SoundVolume,
		IsGameSoundEnabled: &params.IsGameSoundEnabled,
	}, nil
}
