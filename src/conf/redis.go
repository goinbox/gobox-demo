package conf

import "time"

type redisConfJson struct {
	Host                         string `json:"host"`
	Pass                         string `json:"pass"`
	Port                         string `json:"port"`
	LogLevel                     int    `json:"log_level"`
	ConnectTimeoutSeconds        int    `json:"connect_timeout_seconds"`
	RWTimeoutSeconds             int    `json:"rw_timeout_seconds"`
	PoolSize                     int    `json:"pool_size"`
	PoolKeepAliveIntervalSeconds int    `json:"pool_keepalive_interval_seconds"`
	PoolClientMaxIdleSeconds     int    `json:"pool_client_max_idle_seconds"`
}

type RedisConf struct {
	Host                  string
	Pass                  string
	Port                  string
	LogLevel              int
	ConnectTimeout        time.Duration
	RWTimeout             time.Duration
	PoolSize              int
	PoolKeepAliveInterval time.Duration
	PoolClientMaxIdleTime time.Duration
}

var RedisConfList []*RedisConf

func initRedisConf() {
	for _, item := range scJson.Redis {
		RedisConfList = append(RedisConfList,
			&RedisConf{
				Host:                  item.Host,
				Pass:                  item.Pass,
				Port:                  item.Port,
				LogLevel:              item.LogLevel,
				ConnectTimeout:        time.Duration(item.ConnectTimeoutSeconds) * time.Second,
				RWTimeout:             time.Duration(item.RWTimeoutSeconds) * time.Second,
				PoolSize:              item.PoolSize,
				PoolKeepAliveInterval: time.Duration(item.PoolKeepAliveIntervalSeconds) * time.Second,
				PoolClientMaxIdleTime: time.Duration(item.PoolClientMaxIdleSeconds) * time.Second,
			})
	}

}
