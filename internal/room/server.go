package room

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	pb "github.com/quizverse3D/Backend/internal/pb/room"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	room, err := s.svc.CreateRoom(ctx, userUuid, &req.Name, req.Password, &req.MaxPlayers, &req.IsPublic)
	if err != nil {
		log.Printf("failed to create room: %v", err)
		return nil, err
	}

	return &pb.CreateRoomParamsResponse{
		Id:         room.ID.String(),
		Name:       room.Name,
		OwnerId:    room.OwnerUuid.String(),
		OwnerName:  room.OwnerName,
		MaxPlayers: room.MaxPlayers,
		IsPublic:   room.IsPublic,
		CreatedAt:  timestamppb.New(*room.CreatedAt)}, nil
}

func (s *Server) GetRoomById(ctx context.Context, req *pb.GetRoomParamsRequest) (*pb.GetRoomParamsResponse, error) {
	// parse pb
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	room, err := s.svc.GetRoomById(ctx, id, true)
	if err != nil {
		log.Printf("failed to get room info: %v", err)
		return nil, err
	}

	return &pb.GetRoomParamsResponse{
		Id:         room.ID.String(),
		Name:       room.Name,
		OwnerId:    room.OwnerUuid.String(),
		OwnerName:  room.OwnerName,
		MaxPlayers: room.MaxPlayers,
		IsPublic:   room.IsPublic,
		CreatedAt:  timestamppb.New(*room.CreatedAt)}, nil
}
