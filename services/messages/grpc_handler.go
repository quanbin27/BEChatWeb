package messages

import (
	"context"
	"fmt"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/messages"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MessageGrpcHandler struct {
	messageService types.MessageService
	messages.UnimplementedMessageServiceServer
}

func NewGrpcGroupsHandler(grpc *grpc.Server, messageService types.MessageService) {
	grpcHandler := &MessageGrpcHandler{
		messageService: messageService,
	}
	messages.RegisterMessageServiceServer(grpc, grpcHandler)
}
func (h *MessageGrpcHandler) SendMessage(ctx context.Context, req *messages.SendMessageRequest) (*messages.SendMessageResponse, error) {
	messageID, createdAt, err := h.messageService.SendMessage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %v", err)
	}

	return &messages.SendMessageResponse{
		MessageID: messageID,
		CreatedAt: timestamppb.New(createdAt),
	}, nil
}
func (h *MessageGrpcHandler) GetMessages(ctx context.Context, req *messages.GetMessagesRequest) (*messages.GetMessagesResponse, error) {
	ListMessages, err := h.messageService.GetMessages(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %v", err)
	}
	var grpcMessages []*messages.Message
	for _, msg := range ListMessages {
		var messageReplyID int32
		if msg.ReplyMessageID != nil {
			messageReplyID = *msg.ReplyMessageID
		} else {
			messageReplyID = 0
		}

		grpcMessages = append(grpcMessages, &messages.Message{
			ID:             msg.ID,
			UserID:         msg.UserID,
			GroupID:        msg.GroupID,
			Content:        msg.Content,
			MessageReplyID: messageReplyID,
			CreatedAt:      timestamppb.New(msg.CreatedAt),
		})
	}

	return &messages.GetMessagesResponse{
		Messages: grpcMessages,
	}, nil
}

func (h *MessageGrpcHandler) GetLatestMessages(ctx context.Context, req *messages.GetLatestMessagesRequest) (*messages.GetLatestMessagesResponse, error) {
	// Gọi service GetLatestMessages
	message, err := h.messageService.GetLatestMessages(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest message: %v", err)
	}
	var messageReplyID int32
	if message.ReplyMessageID != nil {
		messageReplyID = *message.ReplyMessageID
	} else {
		messageReplyID = 0
	}
	return &messages.GetLatestMessagesResponse{
		Message: &messages.Message{
			ID:             message.ID,
			UserID:         message.UserID,
			GroupID:        message.GroupID,
			Content:        message.Content,
			MessageReplyID: messageReplyID,
			CreatedAt:      timestamppb.New(message.CreatedAt),
		},
	}, nil
}
func (h *MessageGrpcHandler) DeleteMessage(ctx context.Context, req *messages.DeleteMessageRequest) (*messages.DeleteMessageResponse, error) {
	// Gọi service DeleteMessage
	groupID, err := h.messageService.DeleteMessage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete message: %v", err)
	}

	// Trả về phản hồi gRPC
	return &messages.DeleteMessageResponse{
		Status:  "Message deleted successfully",
		GroupID: groupID,
	}, nil
}
