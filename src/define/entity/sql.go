package entity

import "reflect"

const (
	EntityMysqlFieldTag = "mysql"
)

type SqlBaseEntity struct {
	Id       int64  `mysql:"id" json:"id" redis:"id"`
	AddTime  string `mysql:"add_time" json:"add_time" redis:"add_time"`
	EditTime string `mysql:"edit_time" json:"edit_time" redis:"edit_time"`
}

func ReflectSqlColNames(ret reflect.Type) []string {
	var cns []string

	for i := 0; i < ret.NumField(); i++ {
		retf := ret.Field(i)
		if retf.Type.Kind() == reflect.Struct {
			cns = ReflectSqlColNames(retf.Type)
			continue
		}

		if name, ok := retf.Tag.Lookup(EntityMysqlFieldTag); ok {
			cns = append(cns, name)
		}
	}

	return cns
}
