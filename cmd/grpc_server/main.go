package main

import (
	"github.com/quanbin27/gRPC-Web-Chat/pkg/grpc"
	"gorm.io/gorm"
	"log"
)

func main() {
	grpcServer := grpc.NewGRPCServer(":9000")
	grpcServer.Run()
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
