package types

import (
	"context"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/users"
	"gorm.io/gorm"
	"time"
)

type UserStore interface {
	GetUserByID(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) error
}
type UserService interface {
	CreateUser(ctx context.Context, user *users.User) error
}
type User struct {
	ID        int32          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
type RegisterUserPayLoad struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=12"`
}

func FromProto(protoUser *users.User) *User {
	return &User{
		ID:       protoUser.UserID,
		Name:     protoUser.Name,
		Email:    protoUser.Email,
		Password: protoUser.Password,
	}
}
