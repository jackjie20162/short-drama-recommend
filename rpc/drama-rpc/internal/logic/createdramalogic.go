package logic

import (
	"context"

	"short-drama-recommend/rpc/drama-rpc/internal/svc"
	"short-drama-recommend/rpc/drama-rpc/pb/short-drama-recommend/rpc/drama-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDramaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateDramaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDramaLogic {
	return &CreateDramaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateDramaLogic) CreateDrama(in *pb.CreateDramaRequest) (*pb.CreateDramaResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.CreateDramaResponse{}, nil
}
