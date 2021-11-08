package resource

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/pool"

	"gdemo/conf"

	"github.com/goinbox/redis"
)

var redisPool *redis.Pool

func InitRedis(config *conf.RedisConf) {
	redisPool = redis.NewPool(&redis.PConfig{
		Config: pool.Config{
			Size:              config.PoolSize,
			MaxIdleTime:       config.PoolClientMaxIdleTime,
			KeepAliveInterval: config.PoolKeepAliveInterval,
		},
		NewClientFunc: newRedisClientFunc(config),
	})
}

func RedisClient(logger golog.Logger) *redis.Client {
	client, _ := redisPool.Get(logger)

	return client
}

func newRedisClientFunc(config *conf.RedisConf) func() (*redis.Client, error) {
	return func() (*redis.Client, error) {
		rconfig := redis.NewConfig(config.Host, config.Pass, config.Port)
		rconfig.ConnectTimeout = config.ConnectTimeout
		rconfig.ReadTimeout = config.RWTimeout
		rconfig.WriteTimeout = config.RWTimeout

		return redis.NewClient(rconfig, nil), nil
	}
}
