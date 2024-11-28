package groups

import (
	"errors"
	"fmt"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
	"gorm.io/gorm"
)

type GroupStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *GroupStore {
	return &GroupStore{db: db}
}
func (s *GroupStore) CreateGroup(group *types.Group) (int32, error) {
	if err := s.db.Create(group).Error; err != nil {
		return 0, err
	}
	return group.ID, nil
}
func (s *GroupStore) GetGroupByID(id int32) (*types.Group, error) {
	var group types.Group
	result := s.db.Unscoped().Where("id = ?", id).First(&group)
	if result.Error != nil {
		return nil, result.Error
	}
	return &group, nil
}
func (s *GroupStore) DeleteGroup(id int32) error {
	return s.db.Where("id = ?", id).Delete(&types.Group{}).Error
}
func (s *GroupStore) ChangeNameGroup(id int32, name string) error {
	return s.db.Model(&types.Group{}).Where("id = ?", id).Update("name", name).Error
}
func (s *GroupStore) AddMemberToGroup(groupDetail *types.GroupDetail) error {
	return s.db.Create(groupDetail).Error
}
func (s *GroupStore) DeleteGroupDetails(groupID int32, userID int32) error {
	query := s.db.Where("group_id = ?", groupID)
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	result := query.Delete(&types.GroupDetail{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete group details: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("no matching record found")
	}
	return nil
}

func (s *GroupStore) GetRoleIDByUserAndGroup(userID, groupID int32) (int32, error) {
	var groupDetail types.GroupDetail
	err := s.db.Where("user_id = ? AND group_id = ?", userID, groupID).First(&groupDetail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil // Không tìm thấy, tức là không có vai trò
		}
		return 0, err
	}
	return groupDetail.RoleID, nil
}

func (s *GroupStore) UpdateGroup(group *types.Group) error {
	return s.db.Save(group).Error
}
func (s *GroupStore) ChangeAdmin(groupID int32, currentAdminID int32, newAdminID int32) error {
	tx := s.db.Begin()

	if err := tx.Model(&types.GroupDetail{}).
		Where("group_id = ? AND user_id = ? AND role_id = ?", groupID, currentAdminID, 1).
		Update("role_id", 2).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to demote current admin: %v", err)
	}

	if err := tx.Model(&types.GroupDetail{}).
		Where("group_id = ? AND user_id = ?", groupID, newAdminID).
		Update("role_id", 1).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to promote new admin: %v", err)
	}

	// Commit the transaction
	return tx.Commit().Error
}
func (s *GroupStore) GetListMembers(groupID int32) ([]types.Member, error) {
	var members []types.Member
	err := s.db.Table("group_details").
		Select("users.id as user_id, users.name as name, group_details.role_id as role_id").
		Joins("JOIN users ON group_details.user_id = users.id").
		Where("group_details.group_id = ?", groupID).
		Scan(&members).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve members for group %d: %v", groupID, err)
	}
	return members, nil
}

func (s *GroupStore) GetListUserGroups(userID int32) ([]types.Group, error) {
	var groups []types.Group
	err := s.db.Table("groups").
		Select("groups.id, groups.name, groups.member_count").
		Joins("JOIN group_details ON group_details.group_id = groups.id").
		Where("group_details.user_id = ?", userID).
		Scan(&groups).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve groups for user %d: %v", userID, err)
	}
	return groups, nil
}
func (s *GroupStore) GetUserGroupsWithTwoMembers(userID int32) ([]*types.GroupWithMessage, error) {
	var groupsWithMessages []*types.GroupWithMessage

	query := `
		SELECT 
			g.id AS group_id,
			g.name AS group_name,
			(SELECT m.content 
			 FROM messages m 
			 WHERE m.group_id = g.id 
			 ORDER BY m.created_at DESC LIMIT 1) AS latest_message,
			MAX(m.created_at) AS latest_message_time,
			g.member_count,
			(SELECT gd2.user_id 
			 FROM group_details gd2 
			 WHERE gd2.group_id = g.id AND gd2.user_id != ?
			 LIMIT 1) AS other_user_id
		FROM 
			` + "`groups`" + ` g
		JOIN 
			group_details gd ON g.id = gd.group_id
		LEFT JOIN 
			messages m ON g.id = m.group_id
		WHERE 
			gd.user_id = ? AND g.member_count = 2
		GROUP BY 
			g.id, g.name, g.member_count
		ORDER BY 
			latest_message_time DESC
	`

	err := s.db.Raw(query, userID, userID).Scan(&groupsWithMessages).Error
	if err != nil {
		return nil, err
	}

	return groupsWithMessages, nil
}
func (s *GroupStore) GetUserGroupsWithMoreThanTwoMembers(userID int32) ([]*types.GroupWithMessage, error) {
	var groupsWithMessages []*types.GroupWithMessage

	query := `
		SELECT 
			g.id AS group_id,
			g.name AS group_name,
			(SELECT m.content 
			 FROM messages m 
			 WHERE m.group_id = g.id 
			 ORDER BY m.created_at DESC LIMIT 1) AS latest_message,
			MAX(m.created_at) AS latest_message_time,
			g.member_count
		FROM 
			` + "`groups`" + ` g
		JOIN 
			group_details gd ON g.id = gd.group_id
		LEFT JOIN 
			messages m ON g.id = m.group_id
		WHERE 
			gd.user_id = ? AND g.member_count > 2
		GROUP BY 
			g.id, g.name, g.member_count
		ORDER BY 
			latest_message_time DESC
	`

	err := s.db.Raw(query, userID).Scan(&groupsWithMessages).Error
	if err != nil {
		return nil, err
	}

	return groupsWithMessages, nil
}
