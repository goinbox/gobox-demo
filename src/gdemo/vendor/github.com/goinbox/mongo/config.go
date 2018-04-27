package mongo

import (
	"github.com/goinbox/golog"

	"time"
)

const (
	DEFAULT_QUERY_TIMEOUT = 10 * time.Second
)

type Config struct {
	LogLevel int

	Host   string
	Port   string
	User   string
	Pass   string
	DBName string

	QueryTimeout time.Duration
}

func NewConfig(host, port, user, pass, dbname string) *Config {
	return &Config{
		LogLevel: golog.LEVEL_INFO,

		Host:   host,
		Port:   port,
		User:   user,
		Pass:   pass,
		DBName: dbname,

		QueryTimeout: DEFAULT_QUERY_TIMEOUT,
	}
}
