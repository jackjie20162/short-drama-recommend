package svc

import (
 "database/sql"
 "short-drama-recommend/internal/db"
 "short-drama-recommend/api/drama-admin-api/internal/config"
 "github.com/zeromicro/go-zero/zrpc"
 "short-drama-recommend/rpc/drama-rpc/pb"
 paymentpb "short-drama-recommend/rpc/payment-rpc/pb"
)

type ServiceContext struct { Config config.Config; DB *sql.DB; Drama pb.DramaAdminServiceClient; Payment paymentpb.PaymentServiceClient }
func NewServiceContext(c config.Config)*ServiceContext{
 database,err:=db.OpenMySQL(c.Mysql.DataSource);if err!=nil{panic(err)}
 return &ServiceContext{Config:c,DB:database,Drama:pb.NewDramaAdminServiceClient(zrpc.MustNewClient(c.DramaRpc).Conn()),Payment:paymentpb.NewPaymentServiceClient(zrpc.MustNewClient(c.PaymentRpc).Conn())}
}
