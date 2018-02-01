package conf

import "time"

var MysqlConf struct {
	Host                  string
	User                  string
	Pass                  string
	Port                  string
	Name                  string
	RWTimeout             time.Duration
	PoolSize              int
	PoolClientMaxIdleTime time.Duration
}

func initMysqlConf() {
	MysqlConf.Host = scJson.Mysql.Host
	MysqlConf.User = scJson.Mysql.User
	MysqlConf.Pass = scJson.Mysql.Pass
	MysqlConf.Port = scJson.Mysql.Port
	MysqlConf.Name = scJson.Mysql.Name
	MysqlConf.RWTimeout = time.Duration(scJson.Mysql.RWTimeoutSeconds) * time.Second
	MysqlConf.PoolSize = scJson.Mysql.PoolSize
	MysqlConf.PoolClientMaxIdleTime = time.Duration(scJson.Mysql.PoolClientMaxIdleSeconds) * time.Second
}
