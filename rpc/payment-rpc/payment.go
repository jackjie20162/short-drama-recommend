package main

import (
 "flag"
 "fmt"
 "short-drama-recommend/rpc/payment-rpc/internal/config"
 "short-drama-recommend/rpc/payment-rpc/internal/logic"
 "short-drama-recommend/rpc/payment-rpc/internal/svc"
 "short-drama-recommend/rpc/payment-rpc/pb"
 "github.com/zeromicro/go-zero/core/conf"
 "github.com/zeromicro/go-zero/zrpc"
 "google.golang.org/grpc"
)

func main() {
 configFile:=flag.String("f","etc/payment.yaml","the config file")
 flag.Parse()
 var c config.Config
 conf.MustLoad(*configFile,&c)
 ctx:=svc.NewServiceContext(c)
 server:=zrpc.MustNewServer(c.RpcServerConf,func(g *grpc.Server){
  pb.RegisterPaymentServiceServer(g,logic.NewPaymentServer(ctx))
 })
 defer server.Stop()
 fmt.Printf("Starting payment rpc at %s...\n",c.ListenOn)
 server.Start()
}
