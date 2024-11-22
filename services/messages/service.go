package messages

import "github.com/quanbin27/gRPC-Web-Chat/services/types"

type MessageService struct {
	messageStore types.MessageStore
}

func NewMessageService(messageStore types.MessageStore) *MessageService {
	return &MessageService{messageStore}
}
