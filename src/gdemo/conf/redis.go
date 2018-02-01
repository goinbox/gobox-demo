package conf

import "time"

var RedisConf struct {
	Host                  string
	Pass                  string
	Port                  string
	RWTimeout             time.Duration
	PoolSize              int
	PoolKeepAliveInterval time.Duration
	PoolClientMaxIdleTime time.Duration
}

func initRedisConf() {
	RedisConf.Host = scJson.Redis.Host
	RedisConf.Pass = scJson.Redis.Pass
	RedisConf.Port = scJson.Redis.Port
	RedisConf.RWTimeout = time.Duration(scJson.Redis.RWTimeoutSeconds) * time.Second
	RedisConf.PoolSize = scJson.Redis.PoolSize
	RedisConf.PoolKeepAliveInterval = time.Duration(scJson.Redis.PoolKeepAliveIntervalSeconds) * time.Second
	RedisConf.PoolClientMaxIdleTime = time.Duration(scJson.Redis.PoolClientMaxIdleSeconds) * time.Second
}
