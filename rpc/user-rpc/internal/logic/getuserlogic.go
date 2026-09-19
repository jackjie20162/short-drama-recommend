package logic

import (
	"context"

	"short-drama-recommend/rpc/user-rpc/internal/svc"
	"short-drama-recommend/rpc/user-rpc/pb"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetUserLogic) GetUser(in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return &pb.GetUserResponse{
		User: &pb.User{
			Id:        in.UserId,
			Country:   "US",
			Language:  "en",
			Locale:    "en-US",
			Timezone:  "UTC",
		},
	}, nil
}
