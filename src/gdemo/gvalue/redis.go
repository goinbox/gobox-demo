package gvalue

import (
	"github.com/goinbox/redis"
	"github.com/goinbox/golog"

	"gdemo/conf"
)

var RedisClientPool *redis.Pool

func InitRedis() {
	config := &redis.PConfig{NewClientFunc: NewRedisClient}
	config.Size = conf.RedisConf.PoolSize
	config.MaxIdleTime = conf.RedisConf.PoolClientMaxIdleTime
	config.KeepAliveInterval = conf.RedisConf.PoolKeepAliveInterval

	RedisClientPool = redis.NewPool(config)
}

func NewRedisClient() (*redis.Client, error) {
	config := redis.NewConfig(conf.RedisConf.Host, conf.RedisConf.Port, conf.RedisConf.Pass)
	config.LogLevel = golog.LEVEL_DEBUG
	config.ReadTimeout = conf.RedisConf.RWTimeout
	config.WriteTimeout = conf.RedisConf.RWTimeout

	return redis.NewClient(config, nil), nil
}
