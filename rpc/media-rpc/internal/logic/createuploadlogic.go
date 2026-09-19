package logic

import (
	"context"

	"short-drama-recommend/rpc/media-rpc/internal/svc"
	"short-drama-recommend/rpc/media-rpc/pb/short-drama-recommend/rpc/media-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUploadLogic {
	return &CreateUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUploadLogic) CreateUpload(in *pb.CreateUploadRequest) (*pb.CreateUploadResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.CreateUploadResponse{}, nil
}
