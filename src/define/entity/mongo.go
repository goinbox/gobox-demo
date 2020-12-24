package entity

import (
	"reflect"
	"time"
)

const (
	EntityMongoFieldTag = "bson"
	EntityMongoBase     = "MongoBaseEntity"
)

type MongoBaseEntity struct {
	Id       interface{} `bson:"_id" json:"_id"`
	AddTime  time.Time   `bson:"add_time" json:"add_time"`
	EditTime time.Time   `bson:"edit_time" json:"edit_time"`
}

func ReflectMongoColNames(ret reflect.Type) []string {
	var cns []string

	for i := 0; i < ret.NumField(); i++ {
		retf := ret.Field(i)
		if retf.Type.Kind() == reflect.Struct && retf.Name == EntityMongoBase {
			cns = ReflectMongoColNames(retf.Type)
			continue
		}

		if name, ok := retf.Tag.Lookup(EntityMongoFieldTag); ok {
			cns = append(cns, name)
		}
	}

	return cns
}
