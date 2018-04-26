package conf

import "time"

var MongoConf struct {
	Host                  string
	User                  string
	Pass                  string
	Port                  string
	Name                  string
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
	MongoConf.RWTimeout = time.Duration(scJson.Mongo.RWTimeoutSeconds) * time.Second
	MongoConf.PoolSize = scJson.Mongo.PoolSize
	MongoConf.PoolClientMaxIdleTime = time.Duration(scJson.Mongo.PoolClientMaxIdleSeconds) * time.Second
}
