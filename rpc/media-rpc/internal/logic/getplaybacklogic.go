package logic

import (
	"context"

	"short-drama-recommend/rpc/media-rpc/internal/svc"
	"short-drama-recommend/rpc/media-rpc/pb/short-drama-recommend/rpc/media-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPlaybackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPlaybackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPlaybackLogic {
	return &GetPlaybackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPlaybackLogic) GetPlayback(in *pb.GetPlaybackRequest) (*pb.PlaybackResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.PlaybackResponse{}, nil
}
