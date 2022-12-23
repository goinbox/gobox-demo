package misc

import (
	"reflect"

	"github.com/goinbox/mysql"
)

const (
	MysqlUpdateColumnTag          = "mysql_update_column"
	MysqlUpdateColumnTagValueOmit = "omit"
)

type MysqlUpdateColumnValueConvertFunc func(interface{}) interface{}

var MysqlUpdateColumnValueConvertFuncMap = map[string]MysqlUpdateColumnValueConvertFunc{}

func MakeMysqlUpdateColumns(params interface{}) []*mysql.SqlUpdateColumn {
	var result []*mysql.SqlUpdateColumn

	rev := reflect.ValueOf(params)
	if rev.Kind() == reflect.Ptr {
		rev = rev.Elem()
	}
	ret := rev.Type()
	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Ptr {
			if revf.IsNil() {
				continue
			}
			revf = revf.Elem()
		}
		retf := ret.Field(i)
		value := revf.Interface()
		tagValue, ok := retf.Tag.Lookup(MysqlUpdateColumnTag)
		if ok {
			if tagValue == MysqlUpdateColumnTagValueOmit {
				continue
			}
			f, ok := MysqlUpdateColumnValueConvertFuncMap[tagValue]
			if ok {
				value = f(value)
			}
		}

		result = append(result, &mysql.SqlUpdateColumn{
			Name:  mysql.ColumnNameByField(&retf),
			Value: value,
		})
	}

	return result
}
