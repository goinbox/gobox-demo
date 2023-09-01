package test

import (
	"fmt"

	"github.com/goinbox/golog"
	"github.com/goinbox/mysql"

	"gdemo/conf"
	"gdemo/logic/app"
	"gdemo/logic/idgen"
	"gdemo/pcontext"
	"gdemo/resource"
	"gdemo/validate"
)

func InitTestResource(prjHome string) {
	err := conf.Init(prjHome + "/conf/server")
	if err != nil {
		panic(err)
	}

	resource.AccessLogger = Logger()
	resource.InitRedis(conf.ServerConf.Redis)
	_ = resource.InitMySQL(conf.ServerConf.MySQL)

	_ = app.InitAppsCacheData(conf.ServerConf.Misc.AppsDataDir)

	err = validate.Init()
	if err != nil {
		panic(fmt.Errorf("validate.InitTestServer error: %w", err))
	}

	err = idgen.Init("snowflake")
	if err != nil {
		panic(fmt.Errorf("idgen.Init error: %w", err))
	}
}

func Logger() golog.Logger {
	w, _ := golog.NewFileWriter("/dev/stdout", 0)

	return golog.NewSimpleLogger(w, golog.NewSimpleFormater()).EnableColor()
}

func MysqlClient() *mysql.Client {
	return resource.MySQLClient()
}

func Context() *pcontext.Context {
	return pcontext.NewContext(Logger()).SetTraceID("test-trace-id")
}
