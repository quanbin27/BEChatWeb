package messages

import (
	"fmt"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
	"gorm.io/gorm"
	"time"
)

type MessageStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *MessageStore {
	return &MessageStore{db}
}
func (s *MessageStore) SendMessage(msg *types.Message) (int32, time.Time, error) {
	if err := s.db.Create(msg).Error; err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to send message: %v", err)
	}
	return msg.ID, msg.CreatedAt, nil
}
func (s *MessageStore) GetMessages(groupID int32) ([]types.Message, error) {
	var messages []types.Message
	err := s.db.Where("group_id = ?", groupID).Order("created_at").Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get messages for group %d: %v", groupID, err)
	}
	return messages, nil
}
func (s *MessageStore) GetLatestMessages(groupID int32) (types.Message, error) {
	var message types.Message
	err := s.db.Where("group_id = ?", groupID).Order("created_at desc").First(&message).Error
	if err != nil {
		return types.Message{}, fmt.Errorf("failed to get latest message for group %d: %v", groupID, err)
	}
	return message, nil
}
func (s *MessageStore) DeleteMessage(msg *types.Message) (int32, error) {
	// Xóa tin nhắn
	groupID := msg.GroupID
	if err := s.db.Delete(msg).Error; err != nil {
		return -1, fmt.Errorf("failed to delete message: %v", err)
	}
	return groupID, nil
}
func (s *MessageStore) GetMessageByID(messageID int32, msg *types.Message) error {
	if err := s.db.First(msg, messageID).Error; err != nil {
		return fmt.Errorf("failed to find message: %v", err)
	}
	return nil
}
