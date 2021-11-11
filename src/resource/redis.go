package resource

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/redis"

	"gdemo/conf"
)

const (
	dbRedisKey = "main"
)

func InitRedis(config *conf.RedisConf) {
	rconfig := redis.NewConfig(config.Host, config.Pass, config.Port)
	rconfig.DialTimeout = config.ConnectTimeout
	rconfig.ReadTimeout = config.RWTimeout
	rconfig.WriteTimeout = config.RWTimeout

	redis.RegisterDB(dbRedisKey, rconfig)
}

func RedisClient(logger golog.Logger) *redis.Client {
	client, _ := redis.NewClientFromPool(dbRedisKey, logger)

	return client
}
