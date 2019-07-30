package config

import (
	"github.com/go-redis/redis"
	"github.com/domahidizoltan/playground-workflow-engine/trend-analyser/internal/trendanalyser"
)

type Context struct {
	TrendAnalyser trendanalyser.TrendAnalyser
}

func Bootstrap(conf Config) Context {
	redis := connectRedis(conf.Redis)
	repo := trendanalyser.NewSymbolHistoryRepository(redis)

	context := Context {
		TrendAnalyser: trendanalyser.NewTrendAnalyser(repo, AppConfig.App.Threshold),
	}

	return context
}

func connectRedis(config Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: "",
		DB:       0,
	})
}