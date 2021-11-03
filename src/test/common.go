package test

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/mysql"
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
