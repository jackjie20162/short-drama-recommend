package logic

import (
	"context"

	"short-drama-recommend/rpc/drama-rpc/internal/svc"
	"short-drama-recommend/rpc/drama-rpc/pb/short-drama-recommend/rpc/drama-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListEpisodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListEpisodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEpisodeLogic {
	return &ListEpisodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListEpisodeLogic) ListEpisode(in *pb.ListEpisodeRequest) (*pb.ListEpisodeResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.ListEpisodeResponse{}, nil
}
