package logic

import (
	"context"

	"short-drama-recommend/rpc/media-rpc/internal/svc"
	"short-drama-recommend/rpc/media-rpc/pb/short-drama-recommend/rpc/media-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompleteUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompleteUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteUploadLogic {
	return &CompleteUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CompleteUploadLogic) CompleteUpload(in *pb.CompleteUploadRequest) (*pb.MediaAsset, error) {
	// todo: add your logic here and delete this line

	return &pb.MediaAsset{}, nil
}
