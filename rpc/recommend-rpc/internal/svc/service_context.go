package svc

import (
	"database/sql"

	"short-drama-recommend/internal/db"
	"short-drama-recommend/rpc/recommend-rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *sql.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	d, err := db.OpenMySQL(c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{Config: c, DB: d}
}
