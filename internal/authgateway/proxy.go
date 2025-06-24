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
	Prefix     string
	TargetAddr string
	Call       func(ctx context.Context, grpcConn *grpc.ClientConn, userId string, body []byte) (any, error)
}

func NewUserRoute(targetAddr string) GRPCServiceRoute {
	servicePrefix := "/user/api/v1/"
	return GRPCServiceRoute{
		Prefix:     servicePrefix,
		TargetAddr: targetAddr,
		Call: func(ctx context.Context, conn *grpc.ClientConn, userId string, body []byte) (any, error) {
			path := ctx.Value("requestPath").(string)
			client := pb.NewUserServiceClient(conn)

			switch strings.TrimPrefix(path, servicePrefix) {
			case "me":
				var req pb.GetUserRequest
				if err := json.Unmarshal(body, &req); err != nil {
					return nil, err
				}
				req.UserId = userId
				return client.GetUser(ctx, &req)

			default:
				return nil, errors.New("unsupported method path: " + path)
			}
		},
	}
}

// Универсальный HTTP-хендлер для REST → gRPC
func ProxyHandler(routes []GRPCServiceRoute) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId")
		if userId == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		path := r.URL.Path
		var route *GRPCServiceRoute
		for _, rt := range routes {
			if strings.HasPrefix(path, rt.Prefix) {
				route = &rt
				break
			}
		}
		if route == nil {
			http.Error(w, "no matching gRPC route", http.StatusNotFound)
			return
		}

		body := make([]byte, r.ContentLength)
		r.Body.Read(body)

		conn, err := grpc.Dial(route.TargetAddr, grpc.WithInsecure())
		if err != nil {
			http.Error(w, "gRPC connection error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		ctx := context.WithValue(r.Context(), "requestPath", r.URL.Path)

		resp, err := route.Call(ctx, conn, userId.(string), body)
		if err != nil {
			http.Error(w, "gRPC call error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
