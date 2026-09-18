package logic

import (
	"context"
	"short-drama-recommend/rpc/user-rpc/internal/svc"
	"short-drama-recommend/rpc/user-rpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedUserServiceServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{svcCtx: svcCtx}
}

func (s *UserServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return NewGetUserLogic(ctx, s.svcCtx).GetUser(in)
}

func (s *UserServer) mustEmbedUnimplementedUserServiceServer() {
	// Compatibility method for generated protobuf implementations.
}

var _ = codes.OK
var _ = status.New
