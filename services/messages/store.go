package messages

import "gorm.io/gorm"

type MessageStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *MessageStore {
	return &MessageStore{db}
}
