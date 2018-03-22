package demo

import (
	"gdemo/misc"
	"gdemo/svc"

	"github.com/goinbox/gomisc"

	"gdemo/dao"
	"testing"
	"time"
)

type testQueryParamsStruct struct {
	Addtime []string `mysql:"add_time"`
	Name    string   `mysql:"name"`
	Status  int      `mysql:"status"`
}

func TestDemoSvc(t *testing.T) {
	demoSvc := NewDemoSvc(nil, misc.MysqlTestClient(), "gdemo", misc.RedisTestClient())

	ids, err := demoSvc.Insert(
		&DemoEntity{Name: "a1", Status: 0},
		&DemoEntity{Name: "a2", Status: 1},
	)
	t.Log(ids, err)

	entity, err := demoSvc.GetById(1)
	t.Log(entity)

	for _, id := range ids {
		entity, err := demoSvc.GetById(id)
		t.Log(entity, err)

		deleted, err := demoSvc.DeleteById(id)
		t.Log(deleted, err)
	}

	baseEntity := svc.SqlBaseEntity{AddTime: time.Now().Format(gomisc.TimeGeneralLayout())}
	demoSvc.UpdateById(1, &DemoEntity{SqlBaseEntity: baseEntity, Name: "aa", Status: 1}, map[string]bool{"add_time": true, "name": true})

	entities, err := demoSvc.ListByIds(1, 10, 11)
	for _, entity := range entities {
		t.Log("listByIds", entity, err)
	}

	sqp := &svc.SqlQueryParams{
		ParamsStructPtr: &testQueryParamsStruct{
			Addtime: []string{"2017-10-01 17:06:30", "2017-10-30 17:06:30"},
			Name:    "%a%",
			Status:  1,
		},
		Exists: map[string]bool{"add_time": true, "name": true, "status": false},
		Conditions: map[string]string{
			"add_time": dao.SQL_COND_BETWEEN,
			"name":     dao.SQL_COND_LIKE,
			"status":   dao.SQL_COND_EQUAL,
		},
		OrderBy: "id desc", Offset: 0, Cnt: 10,
	}
	entities, err = demoSvc.SimpleQueryAnd(sqp)
	for _, entity := range entities {
		t.Log("listByIds", entity, err)
	}

	total, err := demoSvc.TotalRows(demoSvc.EntityName, DEF_DEMO_TOTAL_ROWS_CACHE_EXPIRE_SECONDS)
	t.Log("total:", total, err)
}
