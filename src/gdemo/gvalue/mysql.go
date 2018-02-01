package gvalue

import (
	"github.com/goinbox/mysql"
	"github.com/goinbox/golog"

	"gdemo/conf"
)

var MysqlClientPool *mysql.Pool

func InitMysql() {
	config := &mysql.PConfig{NewClientFunc: NewMysqlClient}
	config.Size = conf.MysqlConf.PoolSize
	config.MaxIdleTime = conf.MysqlConf.PoolClientMaxIdleTime

	MysqlClientPool = mysql.NewPool(config)
}

func NewMysqlClient() (*mysql.Client, error) {
	config := mysql.NewConfig(conf.MysqlConf.User, conf.MysqlConf.Pass, conf.MysqlConf.Host, conf.MysqlConf.Port, conf.MysqlConf.Name)
	config.LogLevel = golog.LEVEL_DEBUG
	config.ReadTimeout = conf.MysqlConf.RWTimeout
	config.WriteTimeout = conf.MysqlConf.RWTimeout

	return mysql.NewClient(config, nil)
}
