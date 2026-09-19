package logic

import (
	"context"

	"short-drama-recommend/rpc/drama-rpc/internal/svc"
	"short-drama-recommend/rpc/drama-rpc/pb/short-drama-recommend/rpc/drama-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListDramaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminListDramaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListDramaLogic {
	return &AdminListDramaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminListDramaLogic) AdminListDrama(in *pb.AdminListDramaRequest) (*pb.AdminListDramaResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.AdminListDramaResponse{}, nil
}
