package service

import (
	"context"
	"errors"
	"github.com/quanbin27/gRPC-Web-Chat/services/auth"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/users"
	"github.com/quanbin27/gRPC-Web-Chat/services/users/types"
)

type UserService struct {
	userStore types.UserStore
}

func NewUserService(userStore types.UserStore) *UserService {
	return &UserService{userStore: userStore}
}
func (s *UserService) CreateUser(ctx context.Context, user *users.User) error {
	_, err := s.userStore.GetUserByEmail(user.Email)
	if err == nil {
		return errors.New("User already exists")
	}
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return errors.New("Failed to hash password")
	}
	user.Password = hashedPassword
	dbUser := types.FromProto(user)
	return s.userStore.CreateUser(dbUser)
}
