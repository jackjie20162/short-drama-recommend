package svc

import (
 "database/sql"
 "short-drama-recommend/internal/db"
 "short-drama-recommend/internal/repository"
 "short-drama-recommend/rpc/drama-rpc/internal/config"
)

type ServiceContext struct {
 Config config.Config
 DB *sql.DB
 DramaRepo repository.DramaRepository
 DramaAdminRepo repository.DramaAdminRepository
 EpisodeRepo repository.EpisodeRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
 database,err:=db.OpenMySQL(c.Mysql.DataSource);if err!=nil{panic(err)}
 return &ServiceContext{Config:c,DB:database,DramaRepo:repository.NewMySQLDramaRepository(database),DramaAdminRepo:repository.NewMySQLDramaAdminRepository(database),EpisodeRepo:repository.NewMySQLEpisodeRepository(database)}
}