package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/goinbox/golog"

	"time"
)

const (
	DEFAULT_QUERY_TIMEOUT  = 10 * time.Second
	DEFAULT_SOCKET_TIMEOUT = 60 * time.Second
	DEFAULT_SYNC_TIMEOUT   = 7 * time.Second
	DEFAULT_MODE           = mgo.Eventual
)

type Config struct {
	LogLevel int

	Hosts  []string
	User   string
	Pass   string
	DBName string

	// Sets the maximum amount of time to allow the query to run.
	QueryTimeout time.Duration

	// Sets the amount of time to wait for a non-responding socket
	// to the database before it is forcefully closed.
	// The default timeout is 1 minute.
	SocketTimeout time.Duration

	// Sets the amount of time an operation with this session
	// will wait before returning an error in case a connection to a usable
	// server can't be established. Set it to zero to wait forever.
	// The default value is 7 seconds.
	SyncTimeout time.Duration

	// changes the consistency mode for the session.
	Mode mgo.Mode
}

func NewConfig(hosts []string, user, pass, dbname string) *Config {
	return &Config{
		LogLevel: golog.LEVEL_INFO,

		Hosts:  hosts,
		User:   user,
		Pass:   pass,
		DBName: dbname,

		QueryTimeout:  DEFAULT_QUERY_TIMEOUT,
		SocketTimeout: DEFAULT_SOCKET_TIMEOUT,
		SyncTimeout:   DEFAULT_SYNC_TIMEOUT,
		Mode:          DEFAULT_MODE,
	}
}
