package entity

import (
	"database/sql"
	"reflect"
)

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

func ReflectSqlInsertColValues(rev reflect.Value) []interface{} {
	var colValues []interface{}

	ret := rev.Type()
	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct {
			colValues = ReflectSqlInsertColValues(revf)
			continue
		}

		_, ok := ret.Field(i).Tag.Lookup(EntityMysqlFieldTag)
		if ok {
			colValues = append(colValues, revf.Interface())
		}
	}

	return colValues
}

func ReflectSqlEntityScanDests(rev reflect.Value) []interface{} {
	var dests []interface{}

	ret := rev.Type()
	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct {
			dests = ReflectSqlEntityScanDests(revf)
			continue
		}

		_, ok := ret.Field(i).Tag.Lookup(EntityMysqlFieldTag)
		if ok {
			dests = append(dests, revf.Addr().Interface())
		}
	}

	return dests
}

func ReflectSqlQueryRowsToEntityList(rows *sql.Rows, ret reflect.Type, entityList interface{}) error {
	rlistv := reflect.ValueOf(entityList).Elem()

	for rows.Next() {
		rev := reflect.New(ret)
		dests := ReflectSqlEntityScanDests(rev.Elem())
		err := rows.Scan(dests...)
		if err != nil {
			return err
		}
		rlistv.Set(reflect.Append(rlistv, rev))
	}

	return nil
}
