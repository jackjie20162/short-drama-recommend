// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"

	"short-drama-recommend/api/drama-api/internal/svc"
	"short-drama-recommend/api/drama-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedLogic {
	return &GetFeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFeedLogic) GetFeed(req *types.FeedReq) (resp *types.FeedResp, err error) {
	// todo: add your logic here and delete this line

	return
}
