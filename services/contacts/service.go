package contacts

import (
	"context"
	"errors"
	"fmt"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/contacts"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
	"time"
)

type ContactService struct {
	store      types.ContactStore
	groupStore types.GroupStore
	userStore  types.UserStore
}

func NewContactService(store types.ContactStore, groupStore types.GroupStore, userStore types.UserStore) *ContactService {
	return &ContactService{store: store, groupStore: groupStore, userStore: userStore}
}

func (s *ContactService) AddContact(ctx context.Context, req *contacts.AddContactRequest) error {
	if req.UserId == req.ContactUserId {
		return errors.New("can't add yourself")
	}
	contact := &types.Contact{
		UserID:        req.UserId,
		ContactUserID: req.ContactUserId,
		Status:        "PENDING",
	}

	err := s.store.AddContact(contact)
	if err != nil {
		return fmt.Errorf("failed to add contact: " + err.Error())
	}
	return nil
}

func (s *ContactService) RemoveContact(ctx context.Context, req *contacts.RemoveContactRequest) error {
	err := s.store.RemoveContact(req.UserId, req.ContactUserId)
	if err != nil {
		return errors.New("failed to remove contact: " + err.Error())
	}
	return nil
}

func (s *ContactService) AcceptContact(ctx context.Context, req *contacts.AcceptContactRequest) error {
	if req.UserId == 0 || req.ContactUserId == 0 {
		return errors.New("invalid user IDs")
	}
	err := s.store.AcceptContact(req.UserId, req.ContactUserId)
	if err != nil {
		return errors.New("failed to accept contact: " + err.Error())
	}
	exists, err := s.groupStore.CheckGroupExists(req.UserId, req.ContactUserId)
	if err != nil {
		return errors.New("failed to check group existence: " + err.Error())
	}
	if exists {
		return nil
	}
	userName1, err1 := s.userStore.GetNameByID(req.UserId)
	userName2, err2 := s.userStore.GetNameByID(req.ContactUserId)
	if err1 != nil || err2 != nil {
		return errors.New("failed to retrieve user names")
	}
	group := &types.Group{
		Name:        userName1 + ", " + userName2,
		MemberCount: 2,
		CreatedAt:   time.Now(),
	}

	groupID, err := s.groupStore.CreateGroup(group)
	if err != nil {
		return errors.New("failed to create group: " + err.Error())
	}
	err = s.groupStore.AddMemberToGroup(&types.GroupDetail{
		UserID:  req.UserId,
		GroupID: groupID,
		RoleID:  2,
	})
	err = s.groupStore.AddMemberToGroup(&types.GroupDetail{
		UserID:  req.ContactUserId,
		GroupID: groupID,
		RoleID:  2,
	})
	if err != nil {
		return errors.New("failed to add member to group: " + err.Error())
	}
	return nil

}
func (s *ContactService) GetContactsNotInGroup(ctx context.Context, userID int32, groupID int32) ([]*contacts.Contact, error) {
	// Lấy danh sách các contact không thuộc nhóm
	contactsList, err := s.store.GetContactsNotInGroup(userID, groupID)
	if err != nil {
		return nil, err
	}

	// Tạo một slice chứa thông tin các contact
	var ListContacts []*contacts.Contact
	for _, contact := range contactsList {
		// Xác định contactUserID là người bạn cần lấy thông tin
		var targetUserID int32
		if contact.UserID == userID {
			targetUserID = contact.ContactUserID
		} else {
			targetUserID = contact.UserID
		}

		// Đảm bảo không thêm chính bản thân userID
		if targetUserID == userID {
			continue
		}

		// Lấy thông tin user từ userStore
		user, err := s.userStore.GetUserByID(targetUserID)
		if err != nil {
			return nil, err
		}

		// Chuyển các thông tin từ contact và user thành message gRPC
		ListContacts = append(ListContacts, &contacts.Contact{
			UserId:   targetUserID,
			Username: user.Name,
			Email:    user.Email,
			Avatar:   user.Avatar,
		})
	}

	return ListContacts, nil
}

func (s *ContactService) GetContacts(ctx context.Context, userID int32) ([]*contacts.Contact, error) {
	// Lấy danh sách userIDs đã kết bạn từ bảng Contact
	friendIDs, err := s.store.GetFriendIDs(userID)
	if err != nil {
		return nil, err
	}
	users, err := s.userStore.GetUsersByIDs(friendIDs)
	if err != nil {
		return nil, err
	}
	var ListContacts []*contacts.Contact
	for _, user := range users {
		ListContacts = append(ListContacts, &contacts.Contact{
			UserId:   user.ID,
			Username: user.Name,
			Email:    user.Email,
			Avatar:   user.Avatar,
		})
	}
	return ListContacts, nil
}
func (s *ContactService) GetPendingSentContacts(ctx context.Context, userID int32) ([]*contacts.Contact, error) {
	// Lấy danh sách các liên hệ gửi đi (pending)
	sentContacts, err := s.store.GetPendingSentContacts(userID)
	if err != nil {
		return nil, err
	}

	// Lấy danh sách user_id của các liên hệ đó
	var contactUserIDs []int32
	for _, contact := range sentContacts {
		contactUserIDs = append(contactUserIDs, contact.ContactUserID)
	}

	// Truy xuất thông tin chi tiết của các user_id từ bảng User
	users, err := s.userStore.GetUsersByIDs(contactUserIDs)
	if err != nil {
		return nil, err
	}

	// Chuyển đổi thành danh sách contacts.Contact
	var pendingContacts []*contacts.Contact
	for _, user := range users {
		pendingContacts = append(pendingContacts, &contacts.Contact{
			UserId:   user.ID,
			Username: user.Name,
			Email:    user.Email,
			Avatar:   user.Avatar,
		})
	}

	return pendingContacts, nil
}

func (s *ContactService) GetPendingReceivedContacts(ctx context.Context, userID int32) ([]*contacts.Contact, error) {
	receivedContacts, err := s.store.GetPendingReceivedContacts(userID)
	if err != nil {
		return nil, err
	}

	var contactUserIDs []int32
	for _, contact := range receivedContacts {
		contactUserIDs = append(contactUserIDs, contact.UserID)
	}

	// Truy xuất thông tin chi tiết của các user_id từ bảng User
	users, err := s.userStore.GetUsersByIDs(contactUserIDs)
	if err != nil {
		return nil, err
	}

	// Chuyển đổi thành danh sách contacts.Contact
	var pendingContacts []*contacts.Contact
	for _, user := range users {
		pendingContacts = append(pendingContacts, &contacts.Contact{
			UserId:   user.ID,
			Username: user.Name,
			Email:    user.Email,
			Avatar:   user.Avatar,
		})
	}

	return pendingContacts, nil
}

func (s *ContactService) RejectContact(ctx context.Context, userID, contactUserID int32) error {
	err := s.store.RejectContact(userID, contactUserID)
	if err != nil {
		return err
	}
	return nil
}
