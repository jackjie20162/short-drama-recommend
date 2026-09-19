package logic

import (
	"context"

	"short-drama-recommend/rpc/behavior-rpc/internal/svc"
	"short-drama-recommend/rpc/behavior-rpc/pb/short-drama-recommend/rpc/behavior-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecordEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecordEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordEventLogic {
	return &RecordEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RecordEventLogic) RecordEvent(in *pb.RecordEventRequest) (*pb.RecordEventResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.RecordEventResponse{}, nil
}
