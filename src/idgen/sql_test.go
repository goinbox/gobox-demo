package idgen

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/mysql"

	"testing"
)

func TestDemoDaoRead(t *testing.T) {
	idGenter := NewSqlIdGenter(getTestClient())

	for i := 0; i < 10; i++ {
		id, err := idGenter.GenId("demo")
		t.Log(id, err)
	}
}

func getTestClient() *mysql.Client {
	config := mysql.NewConfig("root", "123", "127.0.0.1", "3306", "gobox-demo")

	logger := golog.NewSimpleLogger(golog.NewConsoleWriter(), new(golog.NoopFormater))
	client, _ := mysql.NewClient(config, logger)

	return client
}
