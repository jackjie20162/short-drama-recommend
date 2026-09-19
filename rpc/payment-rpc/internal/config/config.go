package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	Stripe struct {
		SecretKey string
		WebhookSecret string
	}
	Paypal struct {
		ClientID string
		ClientSecret string
		Env string
		WebhookID string
	}
}
