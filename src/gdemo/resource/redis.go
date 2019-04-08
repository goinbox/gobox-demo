package resource

import (
	"gdemo/conf"

	"github.com/goinbox/golog"
	"github.com/goinbox/redis"
)

var RedisClientPoolList []*redis.Pool

func InitRedis() {
	for _, rconf := range conf.RedisConfList {
		config := &redis.PConfig{NewClientFunc: NewRedisClientFunc(rconf)}
		config.Size = rconf.PoolSize
		config.MaxIdleTime = rconf.PoolClientMaxIdleTime
		config.KeepAliveInterval = rconf.PoolKeepAliveInterval

		RedisClientPoolList = append(RedisClientPoolList, redis.NewPool(config))
	}
}

func NewRedisClientFunc(rconf *conf.RedisConf) func() (*redis.Client, error) {
	return func() (*redis.Client, error) {
		config := redis.NewConfig(rconf.Host, rconf.Port, rconf.Pass)
		config.LogLevel = golog.LEVEL_INFO
		config.ConnectTimeout = rconf.ConnectTimeout
		config.ReadTimeout = rconf.RWTimeout
		config.WriteTimeout = rconf.RWTimeout

		return redis.NewClient(config, nil), nil
	}
}
