package groups

import (
	"github.com/quanbin27/gRPC-Web-Chat/services/common/genproto/groups"
	"github.com/quanbin27/gRPC-Web-Chat/services/types"
	"google.golang.org/grpc"
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
