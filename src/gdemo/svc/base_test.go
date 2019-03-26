package svc

import (
	"gdemo/misc"
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
	misc.InitTestSystem()
}
