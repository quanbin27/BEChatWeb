package types

import (
	"gorm.io/gorm"
	"time"
)

type UserStore interface {
	GetUserByID(id int) (*User, error)
}

type User struct {
	ID        int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
