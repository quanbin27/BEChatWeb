package handler

import (
	"context"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/users"
	"github.com/quanbin27/gRPC-Web-Chat/services/users/types"
	"google.golang.org/grpc"
)

type UsersGrpcHandler struct {
	userService types.UserService
	users.UnimplementedUserServiceServer
}

func NewGrpcUsersHandler(grpc *grpc.Server, userService types.UserService) {
	grpcHandler := &UsersGrpcHandler{
		userService: userService,
	}
	users.RegisterUserServiceServer(grpc, grpcHandler)
}
func (h *UsersGrpcHandler) Register(ctx context.Context, req *users.RegisterRequest) (*users.RegisterResponse, error) {
	user := users.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}
	err := h.userService.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}
	res := &users.RegisterResponse{
		Status: "success",
	}
	return res, nil
}
