package types

import (
	"context"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/users"
	"gorm.io/gorm"
	"time"
)

type UserStore interface {
	GetUserByID(id int32) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) error
	UpdateInfo(userID int32, updatedData map[string]interface{}) error
	UpdatePassword(userID int32, password string) error
}
type UserService interface {
	CreateUser(ctx context.Context, user *users.RegisterRequest) error
	CreateJWT(ctx context.Context, login *users.LoginRequest) (string, error)
	UpdateUser(ctx context.Context, update *users.ChangeInfoRequest) error
	UpdatePassword(ctx context.Context, update *users.ChangePasswordRequest) error
}
type User struct {
	ID        int32          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Bio       string         `json:"bio"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
type RegisterUserPayLoad struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=12"`
}
type LoginUserPayLoad struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ChangeInfoPayLoad struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Bio   string `json:"bio" validate:"required"`
}
type ChangePasswordPayLoad struct {
	OldPassword        string `json:"old_password" validate:"required,min=3,max=12"`
	NewPassword        string `json:"new_password" validate:"required,min=3,max=12"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required,min=3,max=12"`
}
