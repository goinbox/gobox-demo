package idgen

import (
	"testing"

	"github.com/goinbox/golog"
	"github.com/goinbox/mongo"
)

func TestMongoDemoDaoRead(t *testing.T) {
	idGenter := NewMongoIdGenter(getMongoTestClient())

	for i := 0; i < 10; i++ {
		id, err := idGenter.GenId("demo")
		t.Log(id, err)
	}
}

func getMongoTestClient() *mongo.Client {
	logger := golog.NewSimpleLogger(golog.NewConsoleWriter(), new(golog.NoopFormater))

	config := mongo.NewConfig([]string{"127.0.0.1:27017"}, "", "", "mydb")
	return mongo.NewClient(config, logger)
}
