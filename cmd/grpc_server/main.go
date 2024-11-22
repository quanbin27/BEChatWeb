package main

import (
	"github.com/quanbin27/gRPC-Web-Chat/config"
	"github.com/quanbin27/gRPC-Web-Chat/db"
	"github.com/quanbin27/gRPC-Web-Chat/pkg/grpc"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
	"gorm.io/gorm"
	"log"
)

func main() {
	dsn := config.Envs.DSN
	db, err := db.NewMySQLStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)
	db.AutoMigrate(types.User{})
	grpcServer := grpc.NewGRPCServer(":9000", db)
	if err := grpcServer.Run(); err != nil {
		log.Fatal(err)
	}
}
func initStorage(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to database")
}
