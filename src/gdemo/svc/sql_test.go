package svc

import (
	"gdemo/resource"

	"github.com/goinbox/mysql"

	"testing"
)

func TestSqlSvcInsertGetListUpdateDelete(t *testing.T) {
	ss := NewSqlSvc([]byte("TestSqlSvcInsertGetListUpdateDelete"), resource.MysqlClientPool, true)

	resource.TestLogger.Notice([]byte("test Insert"))

	item := &demoEntity{
		Name:   "tdj",
		Status: 1,
	}

	tableName, entityName := "demo", "demo"
	ids, err := ss.Insert(tableName, entityName, demoColNames, item)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ids)

	resource.TestLogger.Notice([]byte("test List"))

	var data []*demoEntity
	err = ss.ListByIds(tableName, ids, "add_time desc", demoEntityType, &data)
	if err != nil {
		t.Fatal(err)
	}
	for i, item := range data {
		t.Log(i, item)
	}

	sqp := &SqlQueryParams{
		ParamsStructPtr: &demoEntity{
			Status: 1,
		},
		Exists:     map[string]bool{"status": true},
		Conditions: map[string]string{"status": mysql.SQL_COND_EQUAL},

		OrderBy: "add_time desc",
		Offset:  0,
		Cnt:     10,
	}
	cnt, err := ss.SimpleTotalAnd("demo", sqp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cnt)

	data = []*demoEntity{}
	err = ss.SimpleQueryAnd(tableName, sqp, demoEntityType, &data)
	if err != nil {
		t.Fatal(err)
	}
	for i, item := range data {
		t.Log(i, item)
	}

	newDemo := &demoEntity{
		Name: "new-demo",
	}
	updateFields := map[string]bool{"name": true}
	setItems, err := ss.UpdateById(tableName, ids[0], newDemo, updateFields)
	if err != nil {
		t.Fatal(err)
	}
	for i, item := range setItems {
		t.Log(i, item)
	}

	resource.TestLogger.Notice([]byte("test Get"))

	item = &demoEntity{}
	find, err := ss.GetById(tableName, ids[0], item)
	if !find {
		t.Fatal("not found")
	}
	if err != nil {
		t.Fatal(err)
	}
	t.Log(item)

	resource.TestLogger.Notice([]byte("test Delete"))

	result := ss.Dao().DeleteById(tableName, ids[0])
	defer ss.SendBackClient()

	t.Log(result)
}
