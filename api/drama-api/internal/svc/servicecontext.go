package svc

import (
	"short-drama-recommend/api/drama-api/internal/config"
	behaviorpb "short-drama-recommend/rpc/behavior-rpc/pb"
	dramapb "short-drama-recommend/rpc/drama-rpc/pb"
	recommendpb "short-drama-recommend/rpc/recommend-rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	Drama dramapb.DramaServiceClient
	Behavior behaviorpb.BehaviorServiceClient
	Recommend recommendpb.RecommendServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:c,
		Drama:dramapb.NewDramaServiceClient(zrpc.MustNewClient(c.DramaRpc).Conn()),
		Behavior:behaviorpb.NewBehaviorServiceClient(zrpc.MustNewClient(c.BehaviorRpc).Conn()),
		Recommend:recommendpb.NewRecommendServiceClient(zrpc.MustNewClient(c.RecommendRpc).Conn()),
	}
}
