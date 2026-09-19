package logic

import (
	"context"

	"short-drama-recommend/rpc/drama-rpc/internal/svc"
	"short-drama-recommend/rpc/drama-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDramaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDramaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDramaLogic {
	return &GetDramaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDramaLogic) GetDrama(in *pb.GetDramaRequest) (*pb.GetDramaResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetDramaResponse{}, nil
}
