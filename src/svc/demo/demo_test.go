package demo

import (
	"gdemo/conf"
	"gdemo/resource"
	"gdemo/svc"

	"github.com/goinbox/gomisc"
	"github.com/goinbox/mysql"

	"os"
	"testing"
	"time"
)

type testQueryParamsStruct struct {
	Addtime []string `mysql:"add_time"`
	Name    string   `mysql:"name"`
	Status  int      `mysql:"status"`
}

func init() {
	_ = conf.Init(os.Getenv("GOPATH"))

	_ = resource.InitLog("test")
	resource.InitRedis()
	resource.InitMysql()
}

func TestDemoSvc(t *testing.T) {
	demoSvc := NewDemoSvc([]byte("tracedemosvc"))

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
	_, _ = demoSvc.UpdateById(1, &DemoEntity{SqlBaseEntity: baseEntity, Name: "aa", Status: 1}, map[string]bool{"add_time": true, "name": true})

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
			"add_time": mysql.SQL_COND_BETWEEN,
			"name":     mysql.SQL_COND_LIKE,
			"status":   mysql.SQL_COND_EQUAL,
		},
		OrderBy: "id desc", Offset: 0, Cnt: 10,
	}
	entities, err = demoSvc.SimpleQueryAnd(sqp)
	for _, entity := range entities {
		t.Log("listByIds", entity, err)
	}

	total, merr, rerr := demoSvc.TotalRows(demoSvc.EntityName, demoSvc.RedisKeyPrefix, DefDemoTotalRowsCacheExpireSeconds)
	t.Log("total:", total, merr, rerr)
}
