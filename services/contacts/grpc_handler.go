package contacts

import (
	"context"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/contacts"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ContactsGrpcHandler struct {
	service types.ContactService
	contacts.UnimplementedContactServiceServer
}

func NewGrpcContactsHandler(grpc *grpc.Server, contactService types.ContactService) {
	grpcHandler := &ContactsGrpcHandler{
		service: contactService,
	}
	contacts.RegisterContactServiceServer(grpc, grpcHandler)
}

func (h *ContactsGrpcHandler) AddContact(ctx context.Context, req *contacts.AddContactRequest) (*contacts.AddContactResponse, error) {
	err := h.service.AddContact(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &contacts.AddContactResponse{Status: "SUCCESS"}, nil
}

func (h *ContactsGrpcHandler) RemoveContact(ctx context.Context, req *contacts.RemoveContactRequest) (*contacts.RemoveContactResponse, error) {
	err := h.service.RemoveContact(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &contacts.RemoveContactResponse{Status: "SUCCESS"}, nil
}

func (h *ContactsGrpcHandler) AcceptContact(ctx context.Context, req *contacts.AcceptContactRequest) (*contacts.AcceptContactResponse, error) {
	err := h.service.AcceptContact(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &contacts.AcceptContactResponse{Status: "SUCCESS"}, nil
}

func (h *ContactsGrpcHandler) GetContacts(ctx context.Context, req *contacts.GetContactsRequest) (*contacts.GetContactsResponse, error) {
	ListContacts, err := h.service.GetContacts(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &contacts.GetContactsResponse{Contacts: ListContacts}, nil
}
func (h *ContactsGrpcHandler) GetPendingSentContacts(ctx context.Context, req *contacts.GetPendingSentContactsRequest) (*contacts.GetPendingSentContactsResponse, error) {
	contactsList, err := h.service.GetPendingSentContacts(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	// Chuyển đổi slice của types.Contact thành slice của gRPC contacts.Contact
	var grpcContacts []*contacts.Contact
	for _, contact := range contactsList {
		grpcContacts = append(grpcContacts, &contacts.Contact{
			UserId:        contact.UserID,
			ContactUserId: contact.ContactUserID,
			Status:        contact.Status,
		})
	}

	return &contacts.GetPendingSentContactsResponse{Contacts: grpcContacts}, nil
}

func (h *ContactsGrpcHandler) GetPendingReceivedContacts(ctx context.Context, req *contacts.GetPendingReceivedContactsRequest) (*contacts.GetPendingReceivedContactsResponse, error) {
	contactsList, err := h.service.GetPendingReceivedContacts(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	// Chuyển đổi slice của types.Contact thành slice của gRPC contacts.Contact
	var grpcContacts []*contacts.Contact
	for _, contact := range contactsList {
		grpcContacts = append(grpcContacts, &contacts.Contact{
			UserId:        contact.UserID,
			ContactUserId: contact.ContactUserID,
			Status:        contact.Status,
		})
	}

	return &contacts.GetPendingReceivedContactsResponse{Contacts: grpcContacts}, nil
}

func (h *ContactsGrpcHandler) RejectContact(ctx context.Context, req *contacts.RejectContactRequest) (*contacts.RejectContactResponse, error) {
	err := h.service.RejectContact(ctx, req.UserId, req.ContactUserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &contacts.RejectContactResponse{Status: "SUCCESS"}, nil
}
