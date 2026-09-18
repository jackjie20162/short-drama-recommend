package main
import(
 "flag"
 "fmt"
 "github.com/zeromicro/go-zero/core/conf"
 "github.com/zeromicro/go-zero/rest"
 "github.com/zeromicro/go-zero/zrpc"
 "short-drama-recommend/api/drama-admin-api/internal/config"
 "short-drama-recommend/api/drama-admin-api/internal/handler"
)
func main(){f:=flag.String("f","etc/drama-admin-api.yaml","config file");flag.Parse();var c config.Config;conf.MustLoad(*f,&c);ctx:=handler.NewServiceContext(c);server:=rest.MustNewServer(c.RestConf);defer server.Stop();handler.RegisterRoutes(server,ctx);fmt.Printf("Starting drama admin api at %s\n",c.Host);server.Start();_ = zrpc.MustNewClient}
