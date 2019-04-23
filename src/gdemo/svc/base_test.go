package svc

import (
	"gdemo/conf"
	"gdemo/resource"

	"os"
	"reflect"
)

type demoEntity struct {
	SqlBaseEntity

	Name   string `mysql:"name" json:"name" redis:"name"`
	Status int    `mysql:"status" json:"status" redis:"status"`
}

var demoEntityType reflect.Type = reflect.TypeOf(demoEntity{})
var demoColNames []string = ReflectColNames(demoEntityType)

func init() {
	_ = conf.Init(os.Getenv("GOPATH"))

	_ = resource.InitLog("test")
	resource.InitRedis()
	resource.InitMysql()
	resource.InitMongo()
}
