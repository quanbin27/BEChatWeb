package contacts

import (
	"context"
	"errors"
	"fmt"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/contacts"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
)

type ContactService struct {
	store types.ContactStore
}

func NewContactService(store types.ContactStore) *ContactService {
	return &ContactService{store: store}
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
	err := s.store.AcceptContact(req.UserId, req.ContactUserId)
	if err != nil {
		return errors.New("failed to accept contact: " + err.Error())
	}
	return nil
}

func (s *ContactService) GetContacts(ctx context.Context, req *contacts.GetContactsRequest) ([]*contacts.Contact, error) {
	ListContacts, err := s.store.GetContacts(req.UserId)
	if err != nil {
		return nil, errors.New("failed to fetch contacts: " + err.Error())
	}
	result := make([]*contacts.Contact, len(ListContacts))
	for i, c := range ListContacts {
		result[i] = &contacts.Contact{
			UserId:        c.UserID,
			ContactUserId: c.ContactUserID,
			Status:        c.Status,
		}
	}
	return result, nil
}
func (s *ContactService) GetPendingSentContacts(ctx context.Context, userID int32) ([]types.Contact, error) {
	ListContacts, err := s.store.GetPendingSentContacts(userID)
	if err != nil {
		return nil, err
	}
	return ListContacts, nil
}

func (s *ContactService) GetPendingReceivedContacts(ctx context.Context, userID int32) ([]types.Contact, error) {
	ListContacts, err := s.store.GetPendingReceivedContacts(userID)
	if err != nil {
		return nil, err
	}
	return ListContacts, nil
}
func (s *ContactService) RejectContact(ctx context.Context, userID, contactUserID int32) error {
	err := s.store.RejectContact(userID, contactUserID)
	if err != nil {
		return err
	}
	return nil
}
