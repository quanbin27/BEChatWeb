package groups

import "gorm.io/gorm"

type GroupStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *GroupStore {
	return &GroupStore{db: db}
}
