package logic

import (
	"context"

	"short-drama-recommend/rpc/payment-rpc/internal/svc"
	"short-drama-recommend/rpc/payment-rpc/pb/short-drama-recommend/rpc/payment-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.CreateOrderResponse{}, nil
}
