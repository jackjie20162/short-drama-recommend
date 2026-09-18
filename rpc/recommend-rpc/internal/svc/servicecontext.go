package svc

import (
	"database/sql"

	"short-drama-recommend/internal/cache"
	"short-drama-recommend/internal/db"
	"short-drama-recommend/internal/repository"
	"short-drama-recommend/rpc/recommend-rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB *sql.DB
	DramaRepo repository.DramaRepository
	Redis *cache.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	database, err := db.OpenMySQL(c.Mysql.DataSource)
	if err != nil { panic(err) }
	return &ServiceContext{
		Config:c, DB:database,
		DramaRepo:repository.NewMySQLDramaRepository(database),
		Redis:cache.NewRedis(c.Redis.Host, c.Redis.Pass),
	}
}
