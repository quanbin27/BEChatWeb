package contacts

import (
	"errors"
	"fmt"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
	"gorm.io/gorm"
)

type ContactStore struct {
	db *gorm.DB
}

func NewContactStore(db *gorm.DB) *ContactStore {
	return &ContactStore{db: db}
}

// Thêm liên hệ
func (s *ContactStore) AddContact(contact *types.Contact) error {
	err := s.db.Create(contact).Error
	if err != nil {
		return fmt.Errorf("ContactStore.AddContact: %w", err)
	}
	return nil
}

// Xóa liên hệ
func (s *ContactStore) RemoveContact(userID, contactUserID int32) error {
	// Xóa liên hệ giữa userID và contactUserID, không quan tâm thứ tự.
	res := s.db.Where("user_id = ? AND contact_user_id = ? OR user_id = ? AND contact_user_id = ?",
		userID, contactUserID, contactUserID, userID).
		Delete(&types.Contact{})
	if res.RowsAffected == 0 {
		return errors.New("contact not found")
	}
	return res.Error
}

// Chấp nhận liên hệ
func (s *ContactStore) AcceptContact(userID, contactUserID int32) error {
	res := s.db.Model(&types.Contact{}).
		Where("user_id = ? AND contact_user_id = ? AND status = ?", contactUserID, userID, "PENDING").
		Update("status", "ACCEPTED")
	if res.RowsAffected == 0 {
		return errors.New("contact request not found or already accepted")
	}
	return res.Error
}
func (s *ContactStore) RejectContact(userID, contactUserID int32) error {
	res := s.db.Model(&types.Contact{}).
		Where("user_id = ? AND contact_user_id = ? AND status = ?", contactUserID, userID, "PENDING").
		Update("status", "REJECTED")
	if res.RowsAffected == 0 {
		return errors.New("contact request not found or already accepted")
	}
	return res.Error
}

// Lấy danh sách liên hệ (trạng thái ACCEPTED)
func (s *ContactStore) GetContacts(userID int32) ([]types.Contact, error) {
	var contacts []types.Contact
	err := s.db.Where("(user_id = ? OR contact_user_id = ?) AND status = ?", userID, userID, "ACCEPTED").
		Find(&contacts).Error
	return contacts, err
}
func (s *ContactStore) GetFriendIDs(userID int32) ([]int32, error) {
	var friendIDs []int32
	query := `
		SELECT CASE 
			WHEN user_id = ? THEN contact_user_id 
			ELSE user_id 
		END AS friend_id
		FROM contacts
		WHERE (user_id = ? OR contact_user_id = ?) AND status = 'ACCEPTED'
	`
	err := s.db.Raw(query, userID, userID, userID).Scan(&friendIDs).Error
	if err != nil {
		return nil, err
	}
	return friendIDs, nil
}

func (s *ContactStore) GetPendingSentContacts(userID int32) ([]types.Contact, error) {
	var contacts []types.Contact
	err := s.db.Where("user_id = ? AND status = ?", userID, "PENDING").
		Find(&contacts).Error
	return contacts, err
}
func (s *ContactStore) GetPendingReceivedContacts(userID int32) ([]types.Contact, error) {
	var contacts []types.Contact
	err := s.db.Where("contact_user_id = ? AND status = ?", userID, "PENDING").
		Find(&contacts).Error
	return contacts, err
}
