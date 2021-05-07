package entity

import (
	"reflect"
	"testing"
	"time"

	"github.com/goinbox/gomisc"
)

func TestReflectSqlColNames(t *testing.T) {
	rev := reflect.ValueOf(entityForTest())
	colNames := ReflectSqlColNames(rev.Type().Elem())
	t.Log(colNames)
}

func TestReflectSqlInsertColValues(t *testing.T) {
	rev := reflect.ValueOf(entityForTest())
	colValues := ReflectSqlInsertColValues(rev.Elem())
	t.Log(colValues)
}

func TestReflectSqlEntityScanDests(t *testing.T) {
	rev := reflect.ValueOf(entityForTest())
	dests := ReflectSqlEntityScanDests(rev.Elem())
	t.Log(dests)
}

func entityForTest() *DemoEntity {
	ts := time.Now().Format(gomisc.TimeGeneralLayout())

	return &DemoEntity{
		SqlBaseEntity: SqlBaseEntity{
			Id:       10,
			AddTime:  ts,
			EditTime: ts,
		},
		Name:   "demo",
		Status: 1,
	}
}
