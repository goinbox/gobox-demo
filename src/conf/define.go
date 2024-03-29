package conf

import (
	"time"
)

type LogConf struct {
	Path           string
	Level          int
	AsyncQueueSize int `toml:"async_queue_size"`
	Bufsize        int
	Formater       string
	EnableColor    bool `toml:"enable_color"`
}

type RedisConf struct {
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

type MySQLConf struct {
	Host             string
	User             string
	Pass             string
	Port             int
	Name             string
	RWTimeoutSeconds int `toml:"rw_timeout_seconds"`

	RWTimeout time.Duration
}

type PprofConf struct {
	Enable bool
	Port   int
}

type ApiConf struct {
	Host       string
	Port       int
	EnableSign bool `toml:"enable_sign"`
}

type TracingConf struct {
	Enable           bool
	Target           string
	Async            bool
	ServiceNamespace string `toml:"service_namespace"`

	ZipkinUrl string `toml:"zipkin_url"`
}

var ServerConf struct {
	Hostname string
	Username string

	PrjName string `toml:"prj_name"`
	Idc     string

	Log struct {
		Api *LogConf
	}

	Pprof *PprofConf
	Api   *ApiConf

	Redis *RedisConf
	MySQL *MySQLConf

	Tracing *TracingConf

	Misc struct {
		AppsDataDir string `toml:"apps_data_dir"`
	}
}
