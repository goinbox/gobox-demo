package conf

import "time"

type mysqlConfJson struct {
	Host             string `json:"host"`
	User             string `json:"user"`
	Pass             string `json:"pass"`
	Port             string `json:"port"`
	Name             string `json:"name"`
	LogLevel         int    `json:"log_level"`
	RWTimeoutSeconds int    `json:"rw_timeout_seconds"`
}

var MysqlConf struct {
	Host      string
	User      string
	Pass      string
	Port      string
	Name      string
	LogLevel  int
	RWTimeout time.Duration
}

func initMysqlConf() {
	MysqlConf.Host = scJson.Mysql.Host
	MysqlConf.User = scJson.Mysql.User
	MysqlConf.Pass = scJson.Mysql.Pass
	MysqlConf.Port = scJson.Mysql.Port
	MysqlConf.Name = scJson.Mysql.Name
	MysqlConf.LogLevel = scJson.Mysql.LogLevel
	MysqlConf.RWTimeout = time.Duration(scJson.Mysql.RWTimeoutSeconds) * time.Second
}
