package logic

import (
	"context"

	"short-drama-recommend/rpc/drama-rpc/internal/svc"
	"short-drama-recommend/rpc/drama-rpc/pb/short-drama-recommend/rpc/drama-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetDramaStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetDramaStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDramaStatusLogic {
	return &SetDramaStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetDramaStatusLogic) SetDramaStatus(in *pb.SetDramaStatusRequest) (*pb.SetDramaStatusResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.SetDramaStatusResponse{}, nil
}
