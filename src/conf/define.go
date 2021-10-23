package conf

import (
	"time"
)

type logConf struct {
	Dir            string
	Level          int
	AsyncQueueSize int `json:"async_queue_size"`
	Bufsize        int
}

type redisConf struct {
	Host                         string
	Pass                         string
	Port                         int
	ConnectTimeoutSeconds        int `toml:"connect_timeout_seconds"`
	RWTimeoutSeconds             int `toml:"rw_timeout_seconds"`
	PoolSize                     int `toml:"pool_size"`
	PoolKeepAliveIntervalSeconds int `toml:"pool_keepalive_interval_seconds"`
	PoolClientMaxIdleSeconds     int `toml:"pool_client_max_idle_seconds"`

	ConnectTimeout        time.Duration
	RWTimeout             time.Duration
	PoolKeepAliveInterval time.Duration
	PoolClientMaxIdleTime time.Duration
}

type mysqlConf struct {
	Host             string
	User             string
	Pass             string
	Port             int
	Name             string
	RWTimeoutSeconds int `toml:"rw_timeout_seconds"`

	RWTimeout time.Duration
}

var ServerConf struct {
	Hostname string
	Username string

	PrjName string `toml:"prj_name"`
	IsDev   bool   `toml:"is_dev"`
	Idc     string

	Log   map[string]*logConf
	Pprof struct {
		Enable bool
		Port   int
	}
	Api struct {
		Host    string
		Port    int
		PidFile string `toml:"pid_file"`
	}

	Redis *redisConf
	Mysql *mysqlConf
}
