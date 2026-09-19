package logic

import (
	"context"

	"short-drama-recommend/rpc/drama-rpc/internal/svc"
	"short-drama-recommend/rpc/drama-rpc/pb/short-drama-recommend/rpc/drama-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDramaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDramaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDramaLogic {
	return &UpdateDramaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateDramaLogic) UpdateDrama(in *pb.UpdateDramaRequest) (*pb.UpdateDramaResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateDramaResponse{}, nil
}
