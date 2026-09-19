package main
import("flag";"fmt";"github.com/zeromicro/go-zero/core/conf";"github.com/zeromicro/go-zero/rest";"short-drama-recommend/api/auth-api/internal/config";"short-drama-recommend/api/auth-api/internal/handler")
var configFile=flag.String("f","etc/auth-api.yaml","config")
func main(){flag.Parse();var c config.Config;conf.MustLoad(*configFile,&c);s:=rest.MustNewServer(c.RestConf,rest.WithCors("*"),rest.WithCorsHeaders("Content-Type","Authorization"));defer s.Stop();handler.Register(s,c);fmt.Printf("starting auth api at %s:%d\n",c.Host,c.Port);s.Start()}
