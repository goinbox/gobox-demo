package idgen

import (
	"github.com/goinbox/mysql"
	"github.com/goinbox/golog"

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

	w, _ := golog.NewFileWriter("/tmp/test_mysql.log")
	logger, _ := golog.NewSimpleLogger(w, golog.LEVEL_INFO, golog.NewSimpleFormater())

	client, _ := mysql.NewClient(config, logger)

	return client
}
