package logic

import (
	"context"

	"short-drama-recommend/rpc/drama-rpc/internal/svc"
	"short-drama-recommend/rpc/drama-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDramaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDramaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDramaLogic {
	return &ListDramaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListDramaLogic) ListDrama(in *pb.ListDramaRequest) (*pb.ListDramaResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.ListDramaResponse{}, nil
}
