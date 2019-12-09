package conf

import "time"

type mysqlConfJson struct {
	Host                     string `json:"host"`
	User                     string `json:"user"`
	Pass                     string `json:"pass"`
	Port                     string `json:"port"`
	Name                     string `json:"name"`
	RWTimeoutSeconds         int    `json:"rw_timeout_seconds"`
	PoolSize                 int    `json:"pool_size"`
	PoolClientMaxIdleSeconds int    `json:"pool_client_max_idle_seconds"`
}

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
