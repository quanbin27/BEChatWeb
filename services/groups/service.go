package groups

import "github.com/quanbin27/gRPC-Web-Chat/services/types"

type GroupService struct {
	groupStore types.GroupStore
}

func NewGroupService(groupStore types.GroupStore) *GroupService {
	return &GroupService{groupStore: groupStore}
}
