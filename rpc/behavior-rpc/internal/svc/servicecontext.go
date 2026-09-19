package svc

import (
	"database/sql"

	"short-drama-recommend/internal/cache"
	"short-drama-recommend/internal/db"
	"short-drama-recommend/internal/repository"
	"short-drama-recommend/rpc/behavior-rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB *sql.DB
	BehaviorRepo repository.BehaviorRepository
	Redis *cache.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	database, err := db.OpenMySQL(c.Mysql.DataSource)
	if err != nil { panic(err) }
	r := cache.NewRedis(c.CacheRedis.Host, c.CacheRedis.Pass)
	return &ServiceContext{Config:c, DB:database, BehaviorRepo:repository.NewMySQLBehaviorRepository(database), Redis:r}
}
