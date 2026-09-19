package logic

import (
	"context"

	"short-drama-recommend/rpc/recommend-rpc/internal/svc"
	"short-drama-recommend/rpc/recommend-rpc/pb/short-drama-recommend/rpc/recommend-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedLogic {
	return &GetFeedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFeedLogic) GetFeed(in *pb.FeedRequest) (*pb.FeedResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.FeedResponse{}, nil
}
