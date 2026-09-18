package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/rest"
	"short-drama-recommend/api/drama-api/internal/config"
	"short-drama-recommend/api/drama-api/internal/svc"
)

var configFile = flag.String("f", "etc/drama-api.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	if err := c.Load(*configFile); err != nil { panic(err) }
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	registerRoutes(server, svc.NewServiceContext(c))
	fmt.Printf("starting drama api at %s\n", c.Host+":"+fmt.Sprint(c.Port))
	server.Start()
}
