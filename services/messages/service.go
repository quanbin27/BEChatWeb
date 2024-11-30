package messages

import (
	"context"
	"errors"
	"fmt"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/messages"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
	"time"
)

type MessageService struct {
	store      types.MessageStore
	groupStore types.GroupStore
}

func NewMessageService(messageStore types.MessageStore, groupStore types.GroupStore) *MessageService {
	return &MessageService{messageStore, groupStore}
}
func (s *MessageService) SendMessage(ctx context.Context, req *messages.SendMessageRequest) (int32, time.Time, error) {
	if req.Content == "" {
		return 0, time.Time{}, errors.New("Message content cannot be empty")
	}

	// Lấy ID của tin nhắn reply (nếu có)
	var replyMessageID *int32
	if req.MessageReplyID != 0 {
		replyMessageID = &req.MessageReplyID

		// Kiểm tra xem tin nhắn reply có ở cùng nhóm hay không
		var replyMsg types.Message
		if err := s.store.GetMessageByID(*replyMessageID, &replyMsg); err != nil {
			return 0, time.Time{}, errors.New("Reply message not found or not in the same group")
		}
		if replyMsg.GroupID != req.GroupID {
			return 0, time.Time{}, errors.New("Reply message is not in the same group")
		}
	}

	// Kiểm tra vai trò của người dùng trong nhóm
	roleID, err := s.groupStore.GetRoleIDByUserAndGroup(req.UserID, req.GroupID)
	if err != nil {
		return 0, time.Time{}, err
	}
	if roleID == 0 {
		return 0, time.Time{}, errors.New("User isn't in group")
	}

	// Tạo tin nhắn mới
	msg := &types.Message{
		UserID:         req.UserID,
		GroupID:        req.GroupID,
		Content:        req.Content,
		ReplyMessageID: replyMessageID, // MessageReplyID có thể là nil
		CreatedAt:      time.Now(),
	}
	messageID, createdAt, err := s.store.SendMessage(msg)
	if err != nil {
		return 0, time.Time{}, err
	}

	return messageID, createdAt, nil
}

func (s *MessageService) GetMessages(ctx context.Context, req *messages.GetMessagesRequest) ([]types.Message, error) {
	roleID, err := s.groupStore.GetRoleIDByUserAndGroup(req.UserID, req.GroupID)
	if err != nil {
		return nil, err
	}
	if roleID == 0 {
		return nil, errors.New("User is not authorized to view messages in this group")
	}
	messages, err := s.store.GetMessages(req.GroupID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *MessageService) GetLatestMessages(ctx context.Context, req *messages.GetLatestMessagesRequest) (types.Message, error) {
	roleID, err := s.groupStore.GetRoleIDByUserAndGroup(req.UserID, req.GroupID)
	if err != nil {
		return types.Message{}, err
	}
	if roleID == 0 {
		return types.Message{}, errors.New("User is not authorized to view messages in this group")
	}
	message, err := s.store.GetLatestMessages(req.GroupID)
	if err != nil {
		return types.Message{}, err
	}
	return message, nil
}

func (s *MessageService) DeleteMessage(ctx context.Context, req *messages.DeleteMessageRequest) (int32, error) {
	var msg types.Message
	if err := s.store.GetMessageByID(req.MessageID, &msg); err != nil {
		return -1, fmt.Errorf("failed to find message: %v", err)
	}

	roleID, err := s.groupStore.GetRoleIDByUserAndGroup(req.UserID, msg.GroupID)
	if err != nil {
		return -1, fmt.Errorf("failed to get user role: %v", err)
	}

	if roleID == 1 {
		return s.store.DeleteMessage(&msg)
	}

	if msg.UserID != req.UserID {
		return -1, fmt.Errorf("user %d is not authorized to delete message %d", req.UserID, req.MessageID)
	}

	return s.store.DeleteMessage(&msg)
}
