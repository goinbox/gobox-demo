package test

import (
	"time"

	"github.com/goinbox/golog"
	"github.com/goinbox/mysql"

	"gdemo/conf"
	"gdemo/pcontext"
	"gdemo/resource"
)

func InitMysql() {
	_ = mysql.RegisterDB("test", mysql.NewDefaultConfig("root",
		"123",
		"127.0.0.1",
		"gobox-demo", 3306),
	)
}

func Logger() golog.Logger {
	w, _ := golog.NewFileWriter("/dev/stdout", 0)

	return golog.NewSimpleLogger(w, golog.NewSimpleFormater()).EnableColor()
}

func MysqlClient() *mysql.Client {
	client, _ := mysql.NewClientFromPool("test", Logger())

	return client
}

func InitRedis() {
	resource.InitRedis(&conf.RedisConf{
		Host:                  "127.0.0.1",
		Pass:                  "123",
		Port:                  6379,
		PoolSize:              10,
		ConnectTimeout:        10 * time.Second,
		RWTimeout:             10 * time.Second,
		PoolKeepAliveInterval: 30 * time.Second,
		PoolClientMaxIdleTime: 100 * time.Second,
	})
}

func Context() *pcontext.Context {
	return &pcontext.Context{
		TraceID: "test",
		Logger:  Logger(),
	}
}
