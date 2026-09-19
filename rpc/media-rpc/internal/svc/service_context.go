package svc
import ("database/sql"; "short-drama-recommend/internal/db"; "short-drama-recommend/rpc/media-rpc/internal/config")
type ServiceContext struct { Config config.Config; DB *sql.DB }
func NewServiceContext(c config.Config) *ServiceContext { database,err:=db.OpenMySQL(c.Mysql.DataSource);if err!=nil{panic(err)};return &ServiceContext{Config:c,DB:database} }