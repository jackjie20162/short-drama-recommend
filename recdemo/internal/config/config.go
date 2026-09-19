package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	ES struct {
		Address string
		Index string
	}
	Rank struct {
		PopularityWeight float64
		CompletionWeight float64
		PayWeight float64
		FreshnessWeight float64
		SeenPenalty float64
	}
}
