package groups

import (
	"context"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/groups"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GroupsGrpcHandler struct {
	groupService types.GroupService
	groups.UnimplementedGroupServiceServer
}

func NewGrpcGroupsHandler(grpc *grpc.Server, groupService types.GroupService) {
	grpcHandler := &GroupsGrpcHandler{
		groupService: groupService,
	}
	groups.RegisterGroupServiceServer(grpc, grpcHandler)
}
func (h *GroupsGrpcHandler) CreateGroup(ctx context.Context, req *groups.CreateGroupRequest) (*groups.CreateGroupResponse, error) {
	err := h.groupService.CreateGroup(ctx, req)
	if err != nil {
		return nil, err
	}
	return &groups.CreateGroupResponse{
		Status: "success",
	}, nil
}
func (h *GroupsGrpcHandler) DeleteGroup(ctx context.Context, req *groups.DeleteGroupRequest) (*groups.DeleteGroupResponse, error) {
	err := h.groupService.DeleteGroup(ctx, req)
	if err != nil {
		return nil, err
	}
	return &groups.DeleteGroupResponse{
		Status: "success",
	}, nil
}
func (h *GroupsGrpcHandler) ChangeNameGroup(ctx context.Context, req *groups.ChangeNameRequest) (*groups.ChangeNameResponse, error) {
	newName, err := h.groupService.ChangeNameGroup(ctx, req)
	if err != nil {
		return nil, err
	}
	return &groups.ChangeNameResponse{
		Name:   newName,
		Status: "success",
	}, nil
}
func (h *GroupsGrpcHandler) GetUserGroupsWithLatestMessage(ctx context.Context, req *groups.GetUserGroupsRequest) (*groups.GetUserGroupsResponse, error) {
	groupsList, err := h.groupService.GetUserGroupsWithLatestMessage(ctx, req.UserId, req.MemberCount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get user groups: %v", err)
	}

	var groupResponses []*groups.GroupWithMessage
	for _, group := range groupsList {
		groupResponses = append(groupResponses, &groups.GroupWithMessage{
			GroupId:           group.GroupID,
			GroupName:         group.GroupName,
			LatestMessage:     group.LatestMessage,
			LatestMessageTime: group.LatestMessageTime,
		})
	}

	return &groups.GetUserGroupsResponse{Groups: groupResponses}, nil
}

func (h *GroupsGrpcHandler) GetGroupInfo(ctx context.Context, req *groups.GetGroupInfoRequest) (*groups.Group, error) {
	dbGroup, err := h.groupService.GetGroupInfo(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	res := &groups.Group{
		ID:          dbGroup.ID,
		Name:        dbGroup.Name,
		MemberCount: dbGroup.MemberCount,
	}
	return res, nil
}
func (h *GroupsGrpcHandler) AddMember(ctx context.Context, req *groups.AddMemberRequest) (*groups.AddMemberResponse, error) {
	err := h.groupService.AddMember(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add members: %v", err)
	}
	return &groups.AddMemberResponse{
		Status: "Members added successfully",
	}, nil
}
func (h *GroupsGrpcHandler) KickMember(ctx context.Context, req *groups.KickMemberRequest) (*groups.KickMemberResponse, error) {
	err := h.groupService.KickMember(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to kick member: %v", err)
	}
	return &groups.KickMemberResponse{
		Status: "Member successfully removed from group",
	}, nil
}
func (h *GroupsGrpcHandler) ChangeAdmin(ctx context.Context, req *groups.ChangeAdminRequest) (*groups.ChangeAdminResponse, error) {
	err := h.groupService.ChangeAdmin(ctx, req)
	if err != nil {
		return nil, err
	}

	return &groups.ChangeAdminResponse{
		Status: "Admin changed successfully",
	}, nil
}
func (h *GroupsGrpcHandler) GetListMember(ctx context.Context, req *groups.GetListMemberRequest) (*groups.GetListMemberReponse, error) {
	members, err := h.groupService.GetListMembers(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	var grpcMembers []*groups.Member
	for _, member := range members {
		grpcMembers = append(grpcMembers, &groups.Member{
			UserID: member.ID,
			Name:   member.Name,
			RoleID: member.RoleID,
		})
	}
	response := &groups.GetListMemberReponse{
		Members: grpcMembers,
	}
	return response, nil
}
func (h *GroupsGrpcHandler) LeaveGroup(ctx context.Context, req *groups.LeaveGroupRequest) (*groups.LeaveGroupResponse, error) {
	err := h.groupService.LeaveGroup(ctx, req)
	if err != nil {
		return nil, err
	}

	return &groups.LeaveGroupResponse{
		Status: "Left group successfully",
	}, nil
}
func (h *GroupsGrpcHandler) GetListUserGroup(ctx context.Context, req *groups.GetListUserGroupRequest) (*groups.GetListUserGroupReponse, error) {
	listGroup, err := h.groupService.GetListUserGroups(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	grpcGroups := make([]*groups.Group, len(listGroup))
	for i, g := range listGroup {
		grpcGroups[i] = &groups.Group{
			ID:          g.ID,
			Name:        g.Name,
			MemberCount: g.MemberCount,
		}
	}

	return &groups.GetListUserGroupReponse{
		Groups: grpcGroups,
	}, nil
}
