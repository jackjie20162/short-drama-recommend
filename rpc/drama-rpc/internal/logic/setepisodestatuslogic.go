package logic

import (
	"context"

	"short-drama-recommend/rpc/drama-rpc/internal/svc"
	"short-drama-recommend/rpc/drama-rpc/pb/short-drama-recommend/rpc/drama-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetEpisodeStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetEpisodeStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetEpisodeStatusLogic {
	return &SetEpisodeStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetEpisodeStatusLogic) SetEpisodeStatus(in *pb.SetEpisodeStatusRequest) (*pb.SetEpisodeStatusResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.SetEpisodeStatusResponse{}, nil
}
