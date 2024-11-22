package messages

import (
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/messages"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
	"google.golang.org/grpc"
)

type MessagesGrpcHandler struct {
	messageService types.MessageService
	messages.UnimplementedMessageServiceServer
}

func NewGrpcGroupsHandler(grpc *grpc.Server, messageService types.MessageService) {
	grpcHandler := &MessagesGrpcHandler{
		messageService: messageService,
	}
	messages.RegisterMessageServiceServer(grpc, grpcHandler)
}
