package logic

import (
	"context"

	"short-drama-recommend/rpc/payment-rpc/internal/svc"
	"short-drama-recommend/rpc/payment-rpc/pb/short-drama-recommend/rpc/payment-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrdersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrdersLogic {
	return &ListOrdersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListOrdersLogic) ListOrders(in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.ListOrdersResponse{}, nil
}
