package main

import (
	"flag"
	"fmt"
	"short-drama-recommend/rpc/drama-rpc/internal/config"
	"short-drama-recommend/rpc/drama-rpc/internal/logic"
	"short-drama-recommend/rpc/drama-rpc/internal/svc"
	"short-drama-recommend/rpc/drama-rpc/pb"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	configFile := flag.String("f", "etc/drama.yaml", "the config file")
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	server := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterDramaServiceServer(grpcServer, logic.NewDramaServer(ctx))
	})
	defer server.Stop()
	fmt.Printf("Starting drama rpc at %s...\n", c.ListenOn)
	server.Start()
}
