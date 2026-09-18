package main

import (
	"flag"
	"fmt"

	"short-drama-recommend/rpc/user-rpc/internal/config"
	"short-drama-recommend/rpc/user-rpc/internal/logic"
	"short-drama-recommend/rpc/user-rpc/internal/svc"
	"short-drama-recommend/rpc/user-rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	var configFile = flag.String("f", "etc/user.yaml", "the config file")
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	server := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterUserServiceServer(grpcServer, logic.NewUserServer(ctx))
	})
	defer server.Stop()

	fmt.Printf("Starting user rpc at %s...\n", c.ListenOn)
	server.Start()
}
