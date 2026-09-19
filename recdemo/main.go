package main

import (
 "flag"; "fmt"
 "github.com/zeromicro/go-zero/rest"
 "short-drama-recommend/recdemo/internal/config"
 "short-drama-recommend/recdemo/internal/es"
 "short-drama-recommend/recdemo/internal/handler"
 "short-drama-recommend/recdemo/internal/logic"
 "short-drama-recommend/recdemo/internal/model"
 "short-drama-recommend/recdemo/internal/rank"
)
var configFile=flag.String("f","etc/rec-api.yaml","config file")
func main(){flag.Parse();var c config.Config;if err:=c.Load(*configFile);err!=nil{panic(err)};seeds:=model.Seeds();ranker:=rank.New(c.Rank);esClient:=es.New(es.Config{Address:c.ES.Address,Index:c.ES.Index});mmoe:=rank.NewMMoEClient(c.MMoE.Address,c.MMoE.Model);svc:=logic.NewFeedLogic(seeds,esClient,ranker,mmoe,c.MMoE.Enabled);server:=rest.MustNewServer(c.RestConf);defer server.Stop();handler.Register(server,svc);fmt.Printf("rec-api listening on %s\\n",c.Host);server.Start()}
