package grpc

import (
	"github.com/quanbin27/gRPC-Web-Chat/services/contacts"
	"github.com/quanbin27/gRPC-Web-Chat/services/groups"
	"github.com/quanbin27/gRPC-Web-Chat/services/messages"
	"github.com/quanbin27/gRPC-Web-Chat/services/users"
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
	userService := users.NewUserService(userStore)
	groupStore := groups.NewStore(s.db)
	contactStore := contacts.NewContactStore(s.db)
	groupService := groups.NewGroupService(groupStore)
	messageStore := messages.NewStore(s.db)
	contactService := contacts.NewContactService(contactStore, groupStore, userStore)
	messageService := messages.NewMessageService(messageStore, groupStore)
	users.NewGrpcUsersHandler(grpcServer, userService)
	contacts.NewGrpcContactsHandler(grpcServer, contactService)
	groups.NewGrpcGroupsHandler(grpcServer, groupService)
	messages.NewGrpcGroupsHandler(grpcServer, messageService)
	log.Printf("gRPC pkg listening at %v", s.addr)
	return grpcServer.Serve(lis)
}
