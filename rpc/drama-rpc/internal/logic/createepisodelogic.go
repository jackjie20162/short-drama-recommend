package logic

import (
	"context"

	"short-drama-recommend/rpc/drama-rpc/internal/svc"
	"short-drama-recommend/rpc/drama-rpc/pb/short-drama-recommend/rpc/drama-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateEpisodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateEpisodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEpisodeLogic {
	return &CreateEpisodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateEpisodeLogic) CreateEpisode(in *pb.CreateEpisodeRequest) (*pb.CreateEpisodeResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.CreateEpisodeResponse{}, nil
}
