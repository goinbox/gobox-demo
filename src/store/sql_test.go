package store

import (
	"testing"

	"github.com/goinbox/mysql"

	"gdemo/define"
	"gdemo/define/entity"
	"gdemo/resource"
)

func TestSqlStoreInsertGetListUpdateDelete(t *testing.T) {
	ss := NewSqlStore([]byte("TestSqlStoreInsertGetListUpdateDelete"), resource.MysqlClientPool, true)

	resource.TestLogger.Notice([]byte("test Insert"))

	item := &entity.DemoEntity{
		Name:   "tdj",
		Status: 1,
	}

	tableName, entityName := "demo", "demo"
	ids, err := ss.Insert(tableName, entityName, entity.DemoColNames, item)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ids)

	resource.TestLogger.Notice([]byte("test List"))

	var data []*entity.DemoEntity
	err = ss.ListByIds(tableName, ids, "add_time desc", entity.DemoEntityType, &data)
	if err != nil {
		t.Fatal(err)
	}
	for i, item := range data {
		t.Log(i, item)
	}

	sqp := &define.SqlQueryParams{
		CondItems: []*mysql.SqlColQueryItem{
			{
				Name:      "status",
				Condition: mysql.SqlCondEqual,
				Value:     1,
			},
		},
		OrderBy: "add_time desc",
		Offset:  0,
		Cnt:     10,
	}
	cnt, err := ss.SimpleTotalAnd("demo", sqp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cnt)

	data = []*entity.DemoEntity{}
	err = ss.SimpleQueryAnd(tableName, sqp, entity.DemoEntityType, &data)
	if err != nil {
		t.Fatal(err)
	}
	for i, item := range data {
		t.Log(i, item)
	}

	updateFields := map[string]interface{}{"name": "new-demo"}
	updated, err := ss.UpdateById(tableName, ids[0], updateFields)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(updated)

	resource.TestLogger.Notice([]byte("test Get"))

	item = &entity.DemoEntity{}
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
