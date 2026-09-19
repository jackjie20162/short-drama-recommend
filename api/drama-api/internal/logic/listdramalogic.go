// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"

	"short-drama-recommend/api/drama-api/internal/svc"
	"short-drama-recommend/api/drama-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDramaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListDramaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDramaLogic {
	return &ListDramaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDramaLogic) ListDrama(req *types.ListDramaReq) (resp *types.ListDramaResp, err error) {
	// todo: add your logic here and delete this line

	return
}
