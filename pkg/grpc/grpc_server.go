package grpc

import (
	handler "github.com/quanbin27/gRPC-Web-Chat/services/users/handler/users"
	"github.com/quanbin27/gRPC-Web-Chat/services/users/service"
	users "github.com/quanbin27/gRPC-Web-Chat/services/users/store"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
)

type gRPCServer struct {
	addr string
	db   *gorm.DB
}

func NewGRPCServer(addr string, db *gorm.DB) *gRPCServer {
	return &gRPCServer{addr: addr, db: db}
}
func (s *gRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userStore := users.NewStore(s.db)
	userService := service.NewUserService(userStore)
	handler.NewGrpcUsersHandler(grpcServer, userService)
	log.Printf("gRPC pkg listening at %v", s.addr)
	return grpcServer.Serve(lis)
}
