package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"

	"short-drama-recommend/rpc/behavior-rpc/internal/config"
	"short-drama-recommend/rpc/behavior-rpc/internal/logic"
	"short-drama-recommend/rpc/behavior-rpc/pb"
)

var configFile = flag.String("f", "etc/behavior.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	if err := c.Load(*configFile); err != nil { panic(err) }
	server := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterBehaviorServiceServer(grpcServer, logic.NewBehaviorServer(nil))
	})
	defer server.Stop()
	fmt.Printf("starting behavior rpc at %s\n", c.ListenOn)
	server.Start()
	_ = net.IPv4len
}
