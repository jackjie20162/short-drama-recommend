// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"

	"short-drama-recommend/api/drama-api/internal/svc"
	"short-drama-recommend/api/drama-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecordBehaviorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecordBehaviorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordBehaviorLogic {
	return &RecordBehaviorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecordBehaviorLogic) RecordBehavior(req *types.RecordBehaviorReq) (resp *types.RecordBehaviorResp, err error) {
	// todo: add your logic here and delete this line

	return
}
