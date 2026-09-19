package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"

	"short-drama-recommend/rpc/behavior-rpc/internal/config"
	"short-drama-recommend/rpc/behavior-rpc/internal/logic"
	"short-drama-recommend/rpc/behavior-rpc/internal/svc"
	"short-drama-recommend/rpc/behavior-rpc/pb"
)

var configFile = flag.String("f", "etc/behavior.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)
	server := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterBehaviorServiceServer(grpcServer, logic.NewBehaviorServer(svcCtx))
	})
	defer server.Stop()
	fmt.Printf("starting behavior rpc at %s\n", c.ListenOn)
	server.Start()
}
