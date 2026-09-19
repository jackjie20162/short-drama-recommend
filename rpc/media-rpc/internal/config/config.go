package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	Storage struct {
		Provider string
		Bucket string
		Endpoint string
		Region string
		AccessKey string
		SecretKey string
		PublicBaseURL string
	}
}
