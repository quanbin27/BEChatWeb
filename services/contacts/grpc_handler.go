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
	// Lấy danh sách các liên hệ đã kết bạn
	ListContacts, err := h.service.GetContacts(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get contacts: "+err.Error())
	}

	// Trả về danh sách trong GetContactsResponse
	return &contacts.GetContactsResponse{
		Contacts: ListContacts,
	}, nil
}
func (h *ContactsGrpcHandler) GetPendingSentContacts(ctx context.Context, req *contacts.GetPendingSentContactsRequest) (*contacts.GetPendingSentContactsResponse, error) {
	// Gọi service để lấy danh sách chi tiết thông tin người dùng đã gửi yêu cầu kết bạn
	contactsList, err := h.service.GetPendingSentContacts(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch pending sent contacts: %v", err)
	}

	// Chuyển đổi slice của service response thành slice của gRPC response
	var grpcContacts []*contacts.Contact
	for _, contact := range contactsList {
		grpcContacts = append(grpcContacts, &contacts.Contact{
			UserId:   contact.UserId,
			Username: contact.Username,
			Email:    contact.Email,
		})
	}

	return &contacts.GetPendingSentContactsResponse{Contacts: grpcContacts}, nil
}
func (h *ContactsGrpcHandler) GetPendingReceivedContacts(ctx context.Context, req *contacts.GetPendingReceivedContactsRequest) (*contacts.GetPendingReceivedContactsResponse, error) {
	// Gọi service để lấy danh sách chi tiết thông tin người dùng đã nhận yêu cầu kết bạn
	contactsList, err := h.service.GetPendingReceivedContacts(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch pending received contacts: %v", err)
	}

	// Chuyển đổi slice của service response thành slice của gRPC response
	var grpcContacts []*contacts.Contact
	for _, contact := range contactsList {
		grpcContacts = append(grpcContacts, &contacts.Contact{
			UserId:   contact.UserId,
			Username: contact.Username,
			Email:    contact.Email,
		})
	}

	return &contacts.GetPendingReceivedContactsResponse{Contacts: grpcContacts}, nil
}
func (h *ContactsGrpcHandler) GetContactsNotInGroup(ctx context.Context, req *contacts.GetContactsNotInGroupRequest) (*contacts.GetContactsNotInGroupResponse, error) {
	// Kiểm tra tham số đầu vào
	if req.UserId <= 0 || req.GroupId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user_id or group_id")
	}

	// Gọi service để lấy danh sách contacts không thuộc nhóm
	ListContacts, err := h.service.GetContactsNotInGroup(ctx, req.UserId, req.GroupId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get contacts not in group: %v", err)
	}

	// Trả về response gRPC
	return &contacts.GetContactsNotInGroupResponse{
		Contacts: ListContacts,
	}, nil
}

func (h *ContactsGrpcHandler) RejectContact(ctx context.Context, req *contacts.RejectContactRequest) (*contacts.RejectContactResponse, error) {
	err := h.service.RejectContact(ctx, req.UserId, req.ContactUserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &contacts.RejectContactResponse{Status: "SUCCESS"}, nil
}
