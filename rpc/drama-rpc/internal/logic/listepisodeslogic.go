package logic

import (
	"context"

	"short-drama-recommend/rpc/drama-rpc/internal/svc"
	"short-drama-recommend/rpc/drama-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListEpisodesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListEpisodesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEpisodesLogic {
	return &ListEpisodesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListEpisodesLogic) ListEpisodes(in *pb.ListEpisodesRequest) (*pb.ListEpisodesResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.ListEpisodesResponse{}, nil
}
