package httpclient

import (
	"github.com/goinbox/golog"

	"time"
)

const (
	DEFAULT_TIMEOUT        = 30 * time.Second
	DEFAULT_KEEPALIVE_TIME = 30 * time.Second

	DEFAULT_MAX_IDLE_CONNS_PER_HOST = 10
)

type Config struct {
	LogLevel int

	Timeout       time.Duration
	KeepAliveTime time.Duration

	MaxIdleConnsPerHost int
}

func NewConfig() *Config {
	return &Config{
		LogLevel: golog.LEVEL_INFO,

		Timeout:       DEFAULT_TIMEOUT,
		KeepAliveTime: DEFAULT_KEEPALIVE_TIME,

		MaxIdleConnsPerHost: DEFAULT_MAX_IDLE_CONNS_PER_HOST,
	}
}
