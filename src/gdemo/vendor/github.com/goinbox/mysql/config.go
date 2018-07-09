package mysql

import (
	"github.com/go-sql-driver/mysql"

	"github.com/goinbox/golog"

	"time"
)

const (
	DEFAULT_CONNECT_TIMEOUT = 10 * time.Second
	DEFAULT_READ_TIMEOUT    = 10 * time.Second
	DEFAULT_WRITE_TIMEOUT   = 10 * time.Second
)

type Config struct {
	*mysql.Config

	LogLevel int
}

func NewConfig(user, pass, host, port, dbname string) *Config {
	params := map[string]string{
		"interpolateParams": "true",
	}

	config := &mysql.Config{
		User:                 user,
		Passwd:               pass,
		Net:                  "tcp",
		Addr:                 host + ":" + port,
		DBName:               dbname,
		Params:               params,
		Timeout:              DEFAULT_CONNECT_TIMEOUT,
		ReadTimeout:          DEFAULT_READ_TIMEOUT,
		WriteTimeout:         DEFAULT_WRITE_TIMEOUT,
		AllowNativePasswords: true,
	}

	return &Config{
		Config: config,

		LogLevel: golog.LEVEL_INFO,
	}
}
