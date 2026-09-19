// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"

	"short-drama-recommend/api/drama-api/internal/svc"
	"short-drama-recommend/api/drama-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDramaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDramaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDramaLogic {
	return &GetDramaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDramaLogic) GetDrama(req *types.GetDramaReq) (resp *types.GetDramaResp, err error) {
	// todo: add your logic here and delete this line

	return
}
