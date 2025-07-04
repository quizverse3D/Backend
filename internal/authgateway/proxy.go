package authgateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	pb "github.com/quizverse3D/Backend/internal/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCServiceRoute struct {
	TargetAddr string
	Conn       *grpc.ClientConn
	Call       func(ctx context.Context, grpcConn *grpc.ClientConn, userId string, body []byte) (any, error)
}

func NewUserGrpcServiceRoute(targetAddr string, urlPrefix string) (GRPCServiceRoute, error) {
	// сервис USER
	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return GRPCServiceRoute{}, err
	}
	route := GRPCServiceRoute{
		TargetAddr: targetAddr,
		Conn:       conn,
		Call: func(ctx context.Context, conn *grpc.ClientConn, userId string, body []byte) (any, error) {
			path := strings.TrimPrefix(ctx.Value("requestPath").(string), urlPrefix)
			client := pb.NewUserServiceClient(conn)

			switch path {
			case "me":
				var req pb.GetUserRequest
				if err := json.Unmarshal(body, &req); err != nil {
					return nil, err
				}
				req.UserId = userId
				return client.GetUser(ctx, &req)

			default:
				return nil, errors.New("path not found: " + path)
			}
		},
	}
	return route, nil
}

// Универсальный HTTP-хендлер для REST → gRPC
// Привязка rest-префикса к gRPC-сервису
func ProxyHandler(grpcServiceRoute GRPCServiceRoute) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId")
		if userId == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		ctx := context.WithValue(r.Context(), "requestPath", r.URL.Path)

		resp, err := grpcServiceRoute.Call(ctx, grpcServiceRoute.Conn, userId.(string), body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
