package conf

import "time"

type mongoConfJson struct {
	Host                     string `json:"host"`
	User                     string `json:"user"`
	Pass                     string `json:"pass"`
	Port                     string `json:"port"`
	Name                     string `json:"name"`
	LogLevel                 int    `json:"log_level"`
	RWTimeoutSeconds         int    `json:"rw_timeout_seconds"`
	PoolSize                 int    `json:"pool_size"`
	PoolClientMaxIdleSeconds int    `json:"pool_client_max_idle_seconds"`
}

var MongoConf struct {
	Host                  string
	User                  string
	Pass                  string
	Port                  string
	Name                  string
	LogLevel              int
	RWTimeout             time.Duration
	PoolSize              int
	PoolClientMaxIdleTime time.Duration
}

func initMongoConf() {
	MongoConf.Host = scJson.Mongo.Host
	MongoConf.User = scJson.Mongo.User
	MongoConf.Pass = scJson.Mongo.Pass
	MongoConf.Port = scJson.Mongo.Port
	MongoConf.Name = scJson.Mongo.Name
	MongoConf.LogLevel = scJson.Mongo.LogLevel
	MongoConf.RWTimeout = time.Duration(scJson.Mongo.RWTimeoutSeconds) * time.Second
	MongoConf.PoolSize = scJson.Mongo.PoolSize
	MongoConf.PoolClientMaxIdleTime = time.Duration(scJson.Mongo.PoolClientMaxIdleSeconds) * time.Second
}
