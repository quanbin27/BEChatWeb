package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func NewMySQLStorage() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open("root:27052002@tcp(127.0.0.1:3306)/WebChat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}
