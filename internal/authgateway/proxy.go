package authgateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	pb "github.com/quizverse3D/Backend/internal/pb/user"
	"google.golang.org/grpc"
)

type GRPCServiceRoute struct {
	TargetAddr string
	Call       func(ctx context.Context, grpcConn *grpc.ClientConn, userId string, body []byte) (any, error)
}

func NewUserGrpcServiceRoute(targetAddr string, urlPrefix string) GRPCServiceRoute {
	// сервис USER
	return GRPCServiceRoute{
		TargetAddr: targetAddr,
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

		conn, err := grpc.Dial(grpcServiceRoute.TargetAddr, grpc.WithInsecure())
		if err != nil {
			http.Error(w, "gRPC connection error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		ctx := context.WithValue(r.Context(), "requestPath", r.URL.Path)

		resp, err := grpcServiceRoute.Call(ctx, conn, userId.(string), body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
