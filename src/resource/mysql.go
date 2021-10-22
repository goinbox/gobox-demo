package resource

import (
	"gdemo/conf"

	"github.com/goinbox/mysql"
)

var mysqlClient *mysql.Client

func MysqlClient() (*mysql.Client, error) {


	config := mysql.NewDefaultConfig(conf.MysqlConf.User, conf.MysqlConf.Pass, conf.MysqlConf.Host, conf.MysqlConf.Port, conf.MysqlConf.Name)
	config.LogLevel = conf.MysqlConf.LogLevel
	config.ReadTimeout = conf.MysqlConf.RWTimeout
	config.WriteTimeout = conf.MysqlConf.RWTimeout

	return mysql.NewClient(config, nil)
}
