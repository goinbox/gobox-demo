package mongo

import (
	"github.com/goinbox/golog"

	"time"
)

const (
	DEFAULT_QUERY_TIMEOUT  = 10 * time.Second
	DEFAULT_SOCKET_TIMEOUT = 10 * time.Second
	DEFAULT_SYNC_TIMEOUT   = 7 * time.Second
)

type Config struct {
	LogLevel int

	Hosts   []string
	User   string
	Pass   string
	DBName string

	QueryTimeout  time.Duration
	SocketTimeout time.Duration
	SyncTimeout   time.Duration
}

func NewConfig(hosts []string, user, pass, dbname string) *Config {
	return &Config{
		LogLevel: golog.LEVEL_INFO,

		Hosts:   hosts,
		User:   user,
		Pass:   pass,
		DBName: dbname,

		QueryTimeout:  DEFAULT_QUERY_TIMEOUT,
		SocketTimeout: DEFAULT_SOCKET_TIMEOUT,
		SyncTimeout:   DEFAULT_SYNC_TIMEOUT,
	}
}
