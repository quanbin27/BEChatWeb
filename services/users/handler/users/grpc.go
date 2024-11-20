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
	err := h.userService.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	res := &users.RegisterResponse{
		Status: "success",
	}
	return res, nil
}
func (h *UsersGrpcHandler) Login(ctx context.Context, req *users.LoginRequest) (*users.LoginResponse, error) {
	token, err := h.userService.CreateJWT(ctx, req)
	if err != nil {
		res := &users.LoginResponse{
			Token:  "",
			Status: "fail",
		}
		return res, err
	}
	res := &users.LoginResponse{
		Token:  token,
		Status: "success",
	}
	return res, nil
}
func (h *UsersGrpcHandler) ChangeInfo(ctx context.Context, req *users.ChangeInfoRequest) (*users.ChangeInfoResponse, error) {
	err := h.userService.UpdateUser(ctx, req)
	if err != nil {
		res := &users.ChangeInfoResponse{
			Status: "fail",
		}
		return res, err
	}
	res := &users.ChangeInfoResponse{
		Status: "success",
	}
	return res, nil
}
func (h *UsersGrpcHandler) ChangePassword(ctx context.Context, req *users.ChangePasswordRequest) (*users.ChangePasswordResponse, error) {
	err := h.userService.UpdatePassword(ctx, req)
	if err != nil {
		res := &users.ChangePasswordResponse{
			Status: "fail",
		}
		return res, err
	}
	res := &users.ChangePasswordResponse{
		Status: "success",
	}
	return res, nil
}
