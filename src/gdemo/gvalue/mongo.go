package gvalue

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/mongo"

	"gdemo/conf"
)

var MongoClientPool *mongo.Pool

func InitMongo() {
	config := &mongo.PConfig{NewClientFunc: NewMongoClient}
	config.Size = conf.MongoConf.PoolSize
	config.MaxIdleTime = conf.MongoConf.PoolClientMaxIdleTime

	MongoClientPool = mongo.NewPool(config)
}

func NewMongoClient() (*mongo.Client, error) {
	config := mongo.NewConfig([]string{conf.MongoConf.Host + ":" + conf.MongoConf.Port}, conf.MongoConf.User, conf.MongoConf.Pass, conf.MongoConf.Name)

	config.LogLevel = golog.LEVEL_DEBUG
	config.QueryTimeout = conf.MongoConf.RWTimeout

	return mongo.NewClient(config, nil), nil
}
