package entity

import (
	"reflect"
)

var DemoEntityType = reflect.TypeOf(DemoEntity{})
var DemoColNames = ReflectSqlColNames(DemoEntityType)

type DemoEntity struct {
	SqlBaseEntity

	Name   string `mysql:"name" json:"name" redis:"name"`
	Status int    `mysql:"status" json:"status" redis:"status"`
}

var MongoDemoEntityType = reflect.TypeOf(MongoDemoEntity{})
var MongoDemoColNames = ReflectMongoColNames(MongoDemoEntityType)

type MongoDemoEntity struct {
	MongoBaseEntity

	Name   string `bson:"name" json:"name"`
	Status int    `bson:"status" json:"status"`
}
