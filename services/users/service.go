package users

import (
	"context"
	"errors"
	"github.com/quanbin27/gRPC-Web-Chat/config"
	"github.com/quanbin27/gRPC-Web-Chat/services/auth"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/users"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
)

type UserService struct {
	userStore types.UserStore
}

func NewUserService(userStore types.UserStore) *UserService {
	return &UserService{userStore: userStore}
}
func (s *UserService) CreateUser(ctx context.Context, user *users.RegisterRequest) error {
	_, err := s.userStore.GetUserByEmail(user.Email)
	if err == nil {
		return errors.New("User already exists")
	}
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return errors.New("Failed to hash password")
	}
	return s.userStore.CreateUser(&types.User{Name: user.Name, Email: user.Email, Password: hashedPassword})
}
func (s *UserService) CreateJWT(ctx context.Context, login *users.LoginRequest) (string, error) {
	u, err := s.userStore.GetUserByEmail(login.Email)
	if err != nil {
		return "", errors.New("not found, invalid email")
	}
	if !auth.CheckPassword(u.Password, []byte(login.Password)) {
		return "", errors.New("invalid password")
	}
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID, config.Envs.JWTExpirationInSeconds)
	if err != nil {
		return "", errors.New("Failed to create JWT")
	}
	return token, nil
}
func (s *UserService) UpdateUser(ctx context.Context, update *users.ChangeInfoRequest) error {
	updatedData := map[string]interface{}{
		"name":  update.Name,
		"bio":   update.Bio,
		"email": update.Email,
	}
	err := s.userStore.UpdateInfo(update.Id, updatedData)
	if err != nil {
		return errors.New("Failed to update user")
	}
	return nil
}
func (s *UserService) UpdatePassword(ctx context.Context, update *users.ChangePasswordRequest) error {
	if update.NewPassword == "" {
		return errors.New("Invalid password")
	}
	user, err := s.userStore.GetUserByID(update.Id)
	if err != nil {
		return errors.New("User not found")
	}
	if !auth.CheckPassword(user.Password, []byte(update.OldPassword)) {
		return errors.New("Invalid old password")
	}
	password, err := auth.HashPassword(update.NewPassword)
	err = s.userStore.UpdatePassword(user.ID, password)
	if err != nil {
		return errors.New("Failed to update user")
	}
	return nil
}
func (s *UserService) GetUserByID(ctx context.Context, id int32) (*types.User, error) {
	user, err := s.userStore.GetUserByID(id)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return user, nil
}
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	user, err := s.userStore.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return user, nil
}
