package grpc

import (
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

	log.Printf("gRPC pkg listening at %v", s.addr)
	grpcServer := grpc.NewServer()
	return grpcServer.Serve(lis)
}
