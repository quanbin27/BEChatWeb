package http

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/quanbin27/gRPC-Web-Chat/services/auth"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/contacts"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/groups"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/messages"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/users"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
	"github.com/quanbin27/gRPC-Web-Chat/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"strconv"
	"time"
)

type HttpHandler struct {
	grpcClient *grpc.ClientConn
}

func NewHttpHandler(grpcClient *grpc.ClientConn) *HttpHandler {
	return &HttpHandler{grpcClient: grpcClient}
}
func (h *HttpHandler) RegisterRoutes(e *echo.Group) {
	e.GET("/hello", h.SayHello)
	e.POST("/register", h.RegisterHandler)
	e.POST("/login", h.LoginHandler)
	e.POST("/changeInfo", h.ChangeInfo, auth.WithJWTAuth())
	e.POST("/changePassword", h.ChangePassword, auth.WithJWTAuth())
	e.GET("/user/:id", h.GetUserInfo)
	e.GET("/user/email/:email", h.GetUserByEmail)
	e.GET("/group/:id", h.GetGroupInfo)
	e.POST("/group", h.CreateGroupHandler, auth.WithJWTAuth())
	e.DELETE("/group/:id", h.DeleteGroupHandler, auth.WithJWTAuth())
	e.PATCH("/group/:id/name", h.ChangeNameGroupHandler, auth.WithJWTAuth())
	e.POST("/group/:id/member", h.AddMemberHandler, auth.WithJWTAuth())
	e.DELETE("/group/:id/member", h.KickMemberHandler, auth.WithJWTAuth())
	e.PATCH("/group/:id/change-admin", h.ChangeAdminHandler, auth.WithJWTAuth())
	e.POST("/group/:id/leave", h.LeaveGroupHandler, auth.WithJWTAuth())
	e.GET("/group/:id/member", h.GetListMemberHandler, auth.WithJWTAuth())
	e.GET("/user/group", h.GetListUserGroupsHandler, auth.WithJWTAuth())
	e.GET("/group/:id/message", h.GetMessagesHandler, auth.WithJWTAuth())
	e.POST("/message", h.SendMessageHandler, auth.WithJWTAuth())
	e.DELETE("/message", h.DeleteMessageHandler, auth.WithJWTAuth())
	e.GET("/group/:id/latest-message", h.GetLatestMessagesHandler, auth.WithJWTAuth())
	e.POST("/contact", h.AddContactHandler, auth.WithJWTAuth())
	e.DELETE("/contact", h.RemoveContactHandler, auth.WithJWTAuth())
	e.POST("/contact/accept", h.AcceptContactHandler, auth.WithJWTAuth())
	e.GET("/contacts", h.GetContactsHandler, auth.WithJWTAuth())
	e.GET("/contact/pending-sent", h.GetPendingSentContactsHandler, auth.WithJWTAuth())
	e.GET("/contact/pending-received", h.GetPendingReceivedContactsHandler, auth.WithJWTAuth())
	e.POST("/contact/reject", h.RejectContactHandler, auth.WithJWTAuth())
}
func (h *HttpHandler) GetPendingSentContactsHandler(c echo.Context) error {
	userID := c.Get("user_id").(int32)
	if userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	contactClient := contacts.NewContactServiceClient(h.grpcClient)
	res, err := contactClient.GetPendingSentContacts(ctx, &contacts.GetPendingSentContactsRequest{
		UserId: userID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get pending sent contacts: " + err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
func (h *HttpHandler) GetPendingReceivedContactsHandler(c echo.Context) error {
	userID := c.Get("user_id").(int32)
	if userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	contactClient := contacts.NewContactServiceClient(h.grpcClient)

	res, err := contactClient.GetPendingReceivedContacts(ctx, &contacts.GetPendingReceivedContactsRequest{
		UserId: userID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get pending received contacts: " + err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
func (h *HttpHandler) RejectContactHandler(c echo.Context) error {
	var payload struct {
		ContactUserID int32 `json:"contact_user_id"`
	}
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	userID := c.Get("user_id").(int32)
	if userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	contactClient := contacts.NewContactServiceClient(h.grpcClient)
	res, err := contactClient.RejectContact(ctx, &contacts.RejectContactRequest{
		UserId:        userID,
		ContactUserId: payload.ContactUserID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to reject contact: " + err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
func (h *HttpHandler) AddContactHandler(c echo.Context) error {
	// Bind dữ liệu từ payload request
	var payload types.AddContactPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Kiểm tra tính hợp lệ của payload (ví dụ UserID, ContactUserID)
	if err := utils.Validate.Struct(&payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors[0].Error()})
	}

	// Lấy userID từ context (được thiết lập từ middleware)
	userID := c.Get("user_id").(int32)
	if userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Tạo context và thực hiện yêu cầu
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()

	// Gọi gRPC service để thêm liên hệ
	contactClient := contacts.NewContactServiceClient(h.grpcClient)
	res, err := contactClient.AddContact(ctx, &contacts.AddContactRequest{
		UserId:        userID,
		ContactUserId: payload.ContactUserID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add contact: " + err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
func (h *HttpHandler) RemoveContactHandler(c echo.Context) error {
	// Bind dữ liệu từ payload request
	var payload types.RemoveContactPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Kiểm tra tính hợp lệ của payload (ví dụ UserID, ContactUserID)
	if err := utils.Validate.Struct(&payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors[0].Error()})
	}

	// Lấy userID từ context (được thiết lập từ middleware)
	userID := c.Get("user_id").(int32)
	if userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Tạo context và thực hiện yêu cầu
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()

	// Gọi gRPC service để xóa liên hệ
	contactClient := contacts.NewContactServiceClient(h.grpcClient)
	_, err := contactClient.RemoveContact(ctx, &contacts.RemoveContactRequest{
		UserId:        userID,
		ContactUserId: payload.ContactUserID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove contact: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Contact removed successfully"})
}
func (h *HttpHandler) AcceptContactHandler(c echo.Context) error {
	// Bind dữ liệu từ payload request
	var payload types.AcceptContactPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Kiểm tra tính hợp lệ của payload (ví dụ UserID, ContactUserID)
	if err := utils.Validate.Struct(&payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors[0].Error()})
	}

	// Lấy userID từ context (được thiết lập từ middleware)
	userID := c.Get("user_id").(int32)
	if userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Tạo context và thực hiện yêu cầu
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()

	// Gọi gRPC service để chấp nhận liên hệ
	contactClient := contacts.NewContactServiceClient(h.grpcClient)
	_, err := contactClient.AcceptContact(ctx, &contacts.AcceptContactRequest{
		UserId:        userID,
		ContactUserId: payload.ContactUserID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to accept contact: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Contact accepted successfully"})
}
func (h *HttpHandler) GetContactsHandler(c echo.Context) error {
	// Lấy userID từ context (được thiết lập từ middleware)
	userID := c.Get("user_id").(int32)
	if userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Tạo context và thực hiện yêu cầu
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()

	// Gọi gRPC service để lấy danh sách liên hệ
	contactClient := contacts.NewContactServiceClient(h.grpcClient)
	res, err := contactClient.GetContacts(ctx, &contacts.GetContactsRequest{
		UserId: userID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get contacts: " + err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *HttpHandler) GetLatestMessagesHandler(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil || groupID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid group ID"})
	}
	userID := c.Get("user_id").(int32)
	if userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	messageClient := messages.NewMessageServiceClient(h.grpcClient)
	res, err := messageClient.GetLatestMessages(ctx, &messages.GetLatestMessagesRequest{
		GroupID: int32(groupID),
		UserID:  userID,
	})
	if err != nil {
		if err.Error() == "rpc error: code = PermissionDenied desc = User is not authorized to view messages in this group" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "You do not have permission to access this group"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get latest message: " + err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
func (h *HttpHandler) SendMessageHandler(c echo.Context) error {
	var payload types.SendMessagePayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := utils.Validate.Struct(&payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors[0].Error()})
	}

	userID, ok := c.Get("user_id").(int32)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()

	messageClient := messages.NewMessageServiceClient(h.grpcClient)
	res, err := messageClient.SendMessage(ctx, &messages.SendMessageRequest{
		UserID:         userID,
		GroupID:        payload.GroupID,
		Content:        payload.Content,
		MessageReplyID: getReplyMessageID(payload.MessageReplyID),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send message: " + err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
func (h *HttpHandler) DeleteMessageHandler(c echo.Context) error {
	var payload types.DeleteMessagePayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if err := utils.Validate.Struct(&payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors[0].Error()})
	}

	userID, ok := c.Get("user_id").(int32)
	if !ok || userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()

	messageClient := messages.NewMessageServiceClient(h.grpcClient)
	_, err := messageClient.DeleteMessage(ctx, &messages.DeleteMessageRequest{
		UserID:    userID,
		MessageID: payload.MessageID,
	})
	if err != nil {
		if err.Error() == "rpc error: code = PermissionDenied desc = User is not authorized to delete this message" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "You do not have permission to delete this message"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete message: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Message deleted successfully"})
}

func getReplyMessageID(replyMessageID *int32) int32 {
	if replyMessageID != nil {
		return *replyMessageID
	}
	return 0
}

func (h *HttpHandler) GetMessagesHandler(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil || groupID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid group ID"})
	}

	userID := c.Get("user_id").(int32)
	if userID <= 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()

	// Gọi MessageService thông qua gRPC
	groupClient := messages.NewMessageServiceClient(h.grpcClient)
	res, err := groupClient.GetMessages(ctx, &messages.GetMessagesRequest{
		GroupID: int32(groupID),
		UserID:  userID,
	})
	if err != nil {
		grpcErr := err.Error()
		if grpcErr == "rpc error: code = PermissionDenied desc = User is not authorized to view messages in this group" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "You do not have permission to access this group"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get messages: " + grpcErr})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *HttpHandler) GetListMemberHandler(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil || groupID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid group ID"})
	}

	groupClient := groups.NewGroupServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()

	res, err := groupClient.GetListMember(ctx, &groups.GetListMemberRequest{
		GroupID: int32(groupID),
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
func (h *HttpHandler) GetListUserGroupsHandler(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	groupClient := groups.NewGroupServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()

	res, err := groupClient.GetListUserGroup(ctx, &groups.GetListUserGroupRequest{
		UserID: userID,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *HttpHandler) LeaveGroupHandler(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil || groupID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid group ID"})
	}

	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	groupClient := groups.NewGroupServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()

	_, err = groupClient.LeaveGroup(ctx, &groups.LeaveGroupRequest{
		UserID:  userID,
		GroupID: int32(groupID),
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "Left group successfully"})
}

func (h *HttpHandler) ChangeAdminHandler(c echo.Context) error {
	var payload types.ChangeAdminPayload
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil || groupID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid group ID"})
	}
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	if err := c.Bind(&payload); err != nil || payload.NewAdminID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	groupClient := groups.NewGroupServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	_, err = groupClient.ChangeAdmin(ctx, &groups.ChangeAdminRequest{
		UserID:     userID,
		GroupID:    int32(groupID),
		NewAdminID: payload.NewAdminID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "Admin changed successfully"})
}

func (h *HttpHandler) KickMemberHandler(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil || groupID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid group ID",
		})
	}
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}
	var payload types.KickMemberFromGroupPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	if payload.MemberID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid Member ID",
		})
	}
	groupClient := groups.NewGroupServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()

	res, err := groupClient.KickMember(ctx, &groups.KickMemberRequest{
		UserID:   userID,
		GroupID:  int32(groupID),
		MemberID: payload.MemberID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to kick member: %v", err),
		})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *HttpHandler) AddMemberHandler(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil || groupID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid group ID"})
	}
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	var payload types.AddMemberToGroupPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if len(payload.MemberIds) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Member IDs cannot be empty"})
	}
	groupClient := groups.NewGroupServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	res, err := groupClient.AddMember(ctx, &groups.AddMemberRequest{
		UserID:    userID,
		GroupID:   int32(groupID),
		MemberIDs: payload.MemberIds,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *HttpHandler) GetGroupInfo(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}
	groupClient := groups.NewGroupServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	res, err := groupClient.GetGroupInfo(ctx, &groups.GetGroupInfoRequest{
		GroupID: int32(groupID),
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
func (h *HttpHandler) DeleteGroupHandler(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid group ID")
	}
	groupClient := groups.NewGroupServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	res, err := groupClient.DeleteGroup(ctx, &groups.DeleteGroupRequest{
		GroupID: int32(groupID),
		UserID:  c.Get("user_id").(int32),
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
func (h *HttpHandler) ChangeNameGroupHandler(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid group ID")
	}
	var payload types.ChangeNameGroupPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request name"})
	}
	groupClient := groups.NewGroupServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	res, err := groupClient.ChangeNameGroup(ctx, &groups.ChangeNameRequest{
		GroupID: int32(groupID),
		UserID:  c.Get("user_id").(int32),
		Name:    payload.Name,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
func (h *HttpHandler) CreateGroupHandler(c echo.Context) error {
	var payload types.CreateGroupPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request name"})
	}
	groupClient := groups.NewGroupServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	res, err := groupClient.CreateGroup(ctx, &groups.CreateGroupRequest{
		UserID:    c.Get("user_id").(int32),
		Name:      payload.Name,
		MemberIDs: payload.MemberIds,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
func (h *HttpHandler) ChangeInfo(c echo.Context) error {
	var payload types.ChangeInfoPayLoad
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}
	userClient := users.NewUserServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	res, err := userClient.ChangeInfo(ctx, &users.ChangeInfoRequest{
		Id:    c.Get("user_id").(int32),
		Name:  payload.Name,
		Email: payload.Email,
		Bio:   payload.Bio,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
func (h *HttpHandler) ChangePassword(c echo.Context) error {
	var payload types.ChangePasswordPayLoad
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}
	userClient := users.NewUserServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	res, err := userClient.ChangePassword(ctx, &users.ChangePasswordRequest{
		Id:          c.Get("user_id").(int32),
		OldPassword: payload.OldPassword,
		NewPassword: payload.NewPassword,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
func (h *HttpHandler) LoginHandler(c echo.Context) error {
	var payload types.LoginUserPayLoad
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}
	userClient := users.NewUserServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	res, err := userClient.Login(ctx, &users.LoginRequest{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
func (h *HttpHandler) SayHello(c echo.Context) error {
	return c.JSON(http.StatusOK, "hello world")
}
func (h *HttpHandler) GetUserInfo(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}
	userClient := users.NewUserServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	res, err := userClient.GetUserInfo(ctx, &users.GetUserInfoRequest{
		ID: int32(userID),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)

}
func (h *HttpHandler) GetUserByEmail(c echo.Context) error {
	email := c.Param("email")

	userClient := users.NewUserServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	res, err := userClient.GetUserInfoByEmail(ctx, &users.GetUserInfoByEmailRequest{
		Email: email,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)

}
func (h *HttpHandler) RegisterHandler(c echo.Context) error {
	var payload types.RegisterUserPayLoad
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors[0].Error()})
	}
	userClient := users.NewUserServiceClient(h.grpcClient)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()
	res, err := userClient.Register(ctx, &users.RegisterRequest{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		log.Printf("Failed to register user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func NewGRPCClient(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}
