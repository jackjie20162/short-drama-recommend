package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"

	"short-drama-recommend/rpc/recommend-rpc/internal/config"
	"short-drama-recommend/rpc/recommend-rpc/internal/logic"
	"short-drama-recommend/rpc/recommend-rpc/internal/svc"
	"short-drama-recommend/rpc/recommend-rpc/pb"
)

var configFile = flag.String("f", "etc/recommend.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	if err := c.Load(*configFile); err != nil { panic(err) }
	svcCtx := svc.NewServiceContext(c)
	server := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterRecommendServiceServer(grpcServer, logic.NewRecommendServer(svcCtx))
	})
	defer server.Stop()
	fmt.Printf("starting recommend rpc at %s\n", c.ListenOn)
	server.Start()
}
