package logic

import (
	"context"

	"short-drama-recommend/rpc/payment-rpc/internal/svc"
	"short-drama-recommend/rpc/payment-rpc/pb/short-drama-recommend/rpc/payment-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CapturePaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCapturePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CapturePaymentLogic {
	return &CapturePaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CapturePaymentLogic) CapturePayment(in *pb.CapturePaymentRequest) (*pb.CapturePaymentResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.CapturePaymentResponse{}, nil
}
