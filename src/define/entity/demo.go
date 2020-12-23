package entity

import (
	"reflect"
)

var DemoEntityType reflect.Type = reflect.TypeOf(DemoEntity{})
var DemoColNames []string = ReflectColNames(DemoEntityType)

type DemoEntity struct {
	SqlBaseEntity

	Name   string `mysql:"name" json:"name" redis:"name"`
	Status int    `mysql:"status" json:"status" redis:"status"`
}
