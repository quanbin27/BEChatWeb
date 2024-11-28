package groups

import (
	"context"
	"errors"
	"fmt"
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/groups"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type GroupService struct {
	groupStore   types.GroupStore
	messageStore types.MessageStore
}

func NewGroupService(groupStore types.GroupStore, messageStore types.MessageStore) *GroupService {
	return &GroupService{groupStore: groupStore, messageStore: messageStore}
}
func (s *GroupService) CreateGroup(ctx context.Context, req *groups.CreateGroupRequest) error {
	if req.Name == "" {
		return errors.New("group name is required")
	}

	group := &types.Group{
		Name:        req.Name,
		MemberCount: int32(1 + len(req.MemberIDs)),
		CreatedAt:   time.Now(),
	}
	groupID, err := s.groupStore.CreateGroup(group)
	if err != nil {
		return err
	}
	fisrtMessage := &types.Message{
		UserID:  req.UserID,
		GroupID: groupID,
		Content: "You have been added to the group " + req.Name + " .",
	}
	_, _, err = s.messageStore.SendMessage(fisrtMessage)
	if err != nil {
		return err
	}
	err = s.groupStore.AddMemberToGroup(&types.GroupDetail{
		UserID:  req.UserID,
		GroupID: groupID,
		RoleID:  1,
	})

	if err != nil {
		return err
	}
	for _, memberID := range req.MemberIDs {
		err := s.groupStore.AddMemberToGroup(&types.GroupDetail{
			UserID:  memberID,
			GroupID: groupID,
			RoleID:  2, // Vai trò thành viên
		})
		if err != nil {
			return errors.New("failed to add member %d to group:")
		}
	}
	return nil
}
func (s *GroupService) GetUserGroupsWithLatestMessage(ctx context.Context, userID, memberCount int32) ([]*types.GroupWithMessage, error) {
	// Lấy dữ liệu từ store
	listGroupsWithMessages, err := s.groupStore.GetUserGroupsByFilter(userID, memberCount)
	if err != nil {
		return nil, err
	}
	return listGroupsWithMessages, nil
}

func (s *GroupService) DeleteGroup(ctx context.Context, req *groups.DeleteGroupRequest) error {
	_, err := s.groupStore.GetGroupByID(req.GroupID)
	if err != nil {
		return err
	}
	roleID, err := s.groupStore.GetRoleIDByUserAndGroup(req.UserID, req.GroupID)
	if err != nil {
		return err
	}
	if roleID != 1 {
		return errors.New("You must be admin to delete group")
	}
	err = s.groupStore.DeleteGroupDetails(req.GroupID, 0)
	if err != nil {
		return err
	}
	err = s.groupStore.DeleteGroup(req.GroupID)
	if err != nil {
		return err
	}
	return nil
}
func (s *GroupService) ChangeNameGroup(ctx context.Context, req *groups.ChangeNameRequest) (string, error) {
	_, err := s.groupStore.GetGroupByID(req.GroupID)
	if err != nil {
		return "", err
	}
	roleID, err := s.groupStore.GetRoleIDByUserAndGroup(req.UserID, req.GetGroupID())
	if err != nil {
		return "", err
	}
	if roleID != 1 {
		return "", errors.New("You must be admin to change name group")
	}
	err = s.groupStore.ChangeNameGroup(req.GroupID, req.Name)
	if err != nil {
		return "", err
	}
	return req.Name, nil
}
func (s *GroupService) GetGroupInfo(ctx context.Context, id int32) (*types.Group, error) {
	group, err := s.groupStore.GetGroupByID(id)
	if err != nil {
		return nil, err
	}
	return group, nil
}
func (s *GroupService) AddMember(ctx context.Context, req *groups.AddMemberRequest) error {
	if req.GroupID == 0 || req.UserID == 0 {
		return errors.New("group ID and user ID are required")
	}
	if len(req.MemberIDs) == 0 {
		return errors.New("no members to add")
	}
	roleID, err := s.groupStore.GetRoleIDByUserAndGroup(req.UserID, req.GroupID)
	if err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}
	if roleID != 1 {
		return errors.New("only admin can add members")
	}
	for _, memberID := range req.MemberIDs {
		groupDetail := &types.GroupDetail{
			UserID:    memberID,
			GroupID:   req.GroupID,
			RoleID:    2, // Quyền mặc định là thành viên thường
			CreatedAt: time.Now(),
		}

		if err := s.groupStore.AddMemberToGroup(groupDetail); err != nil {
			return fmt.Errorf("failed to add member with ID %d: %w", memberID, err)
		}
	}
	group, err := s.groupStore.GetGroupByID(req.GroupID)
	if err != nil {
		return err
	}
	group.MemberCount = group.MemberCount + int32(len(req.MemberIDs))
	err = s.groupStore.UpdateGroup(group)
	if err != nil {
		return err
	}
	return nil
}
func (s *GroupService) KickMember(ctx context.Context, req *groups.KickMemberRequest) error {
	if req.UserID == req.MemberID {
		return errors.New("you cannot kick yourself from the group")
	}
	roleID, err := s.groupStore.GetRoleIDByUserAndGroup(req.UserID, req.GroupID)
	if err != nil {
		return fmt.Errorf("failed to get role for user %d: %w", req.UserID, err)
	}
	if roleID != 1 {
		return errors.New("you must be an admin to kick a member")
	}
	memberRoleID, err := s.groupStore.GetRoleIDByUserAndGroup(req.MemberID, req.GroupID)
	if err != nil {
		return fmt.Errorf("failed to get role for member %d: %w", req.MemberID, err)
	}
	if memberRoleID == 0 {
		return errors.New("the member is not part of the group")
	}
	err = s.groupStore.DeleteGroupDetails(req.GroupID, req.MemberID)
	if err != nil {
		return fmt.Errorf("failed to remove member %d from group %d: %w", req.MemberID, req.GroupID, err)
	}
	return nil
}
func (s *GroupService) ChangeAdmin(ctx context.Context, req *groups.ChangeAdminRequest) error {
	roleID, err := s.groupStore.GetRoleIDByUserAndGroup(req.UserID, req.GroupID)
	if err != nil {
		return err
	}
	if roleID != 1 {
		return errors.New("you must be an admin to change the group admin")
	}

	err = s.groupStore.ChangeAdmin(req.GroupID, req.UserID, req.NewAdminID)
	if err != nil {
		return err
	}

	return nil
}
func (s *GroupService) GetListMembers(ctx context.Context, groupID int32) ([]types.Member, error) {
	return s.groupStore.GetListMembers(groupID)
}
func (s *GroupService) LeaveGroup(ctx context.Context, req *groups.LeaveGroupRequest) error {
	roleID, err := s.groupStore.GetRoleIDByUserAndGroup(req.UserID, req.GroupID)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get role: %v", err)
	}
	if roleID == 1 {
		return status.Errorf(codes.FailedPrecondition, "you can't leave the group when you are admin")
	}

	err = s.groupStore.DeleteGroupDetails(req.GroupID, req.UserID)
	if err != nil {
		if err.Error() == "no matching record found" {
			return status.Errorf(codes.NotFound, "record not found for user %d in group %d", req.UserID, req.GroupID)
		}
		return status.Errorf(codes.Internal, "failed to delete group details: %v", err)
	}
	return nil
}

func (s *GroupService) GetListUserGroups(ctx context.Context, userID int32) ([]types.Group, error) {
	return s.groupStore.GetListUserGroups(userID)
}
