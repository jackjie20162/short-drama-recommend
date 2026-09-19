package logic

import (
	"context"

	"short-drama-recommend/rpc/payment-rpc/internal/svc"
	"short-drama-recommend/rpc/payment-rpc/pb/short-drama-recommend/rpc/payment-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleWebhookLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHandleWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleWebhookLogic {
	return &HandleWebhookLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HandleWebhookLogic) HandleWebhook(in *pb.WebhookRequest) (*pb.WebhookResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.WebhookResponse{}, nil
}
