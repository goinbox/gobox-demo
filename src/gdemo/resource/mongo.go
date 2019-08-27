package resource

import (
	"gdemo/conf"

	"github.com/goinbox/golog"
	"github.com/goinbox/mongo"
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

func InitMongoTest() {
	config := &mongo.PConfig{NewClientFunc: NewMongoTestClient}
	config.Size = conf.MongoConf.PoolSize
	config.MaxIdleTime = conf.MongoConf.PoolClientMaxIdleTime

	MongoClientPool = mongo.NewPool(config)
}

func NewMongoTestClient() (*mongo.Client, error) {
	config := mongo.NewConfig([]string{"myhost:myport"}, "myname", "mypass", "mydb")

	config.LogLevel = golog.LEVEL_DEBUG
	config.QueryTimeout = conf.MongoConf.RWTimeout

	return mongo.NewClient(config, nil), nil
}
