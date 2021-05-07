package demo

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gdemo/conf"
	"gdemo/define"
	"gdemo/define/entity"
	"gdemo/resource"

	"github.com/goinbox/gomisc"
	"github.com/goinbox/mysql"
)

type testQueryParamsStruct struct {
	Addtime []string `mysql:"add_time"`
	Name    string   `mysql:"name"`
	Status  int      `mysql:"status"`
}

func init() {
	curDir, _ := os.Getwd()
	prjHome := curDir + "/../../../"
	prjHome, _ = filepath.Abs(prjHome)

	_ = conf.Init(prjHome)

	_ = resource.InitLog("test")
	resource.InitRedis()
	resource.InitMysql()
}

func TestDemoSvc(t *testing.T) {
	demoSvc := NewDemoSvc([]byte("tracedemosvc"))

	ids, err := demoSvc.Insert(
		&entity.DemoEntity{Name: "a1", Status: 0},
		&entity.DemoEntity{Name: "a2", Status: 1},
	)
	t.Log(ids, err)

	item, err := demoSvc.GetById(1)
	t.Log(item)

	for _, id := range ids {
		item, err := demoSvc.GetById(id)
		t.Log(item, err)

		deleted, err := demoSvc.DeleteById(id)
		t.Log(deleted, err)
	}

	updateFields := map[string]interface{}{
		"name":      "aa",
		"status":    1,
		"edit_time": time.Now().Format(gomisc.TimeGeneralLayout()),
	}
	_, _ = demoSvc.UpdateById(1, updateFields)

	entities, err := demoSvc.ListByIds(1, 10, 11)
	for _, item := range entities {
		t.Log("listByIds", item, err)
	}

	sqp := &define.SqlQueryParams{
		CondItems: []*mysql.SqlColQueryItem{
			{
				Name:      "add_time",
				Condition: mysql.SqlCondBetween,
				Value:     []string{"2017-10-01 17:06:30", "2017-10-30 17:06:30"},
			},
			{
				Name:      "name",
				Condition: mysql.SqlCondLike,
				Value:     "%a%",
			},
			{
				Name:      "status",
				Condition: mysql.SqlCondEqual,
				Value:     1,
			},
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
