package demo

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goinbox/mysql"

	"gdemo/model"
	"gdemo/model/demo"
	"gdemo/model/factory"
	"gdemo/model/idgen"
	"gdemo/test"
)

func init() {
	dir, _ := os.Getwd()
	for i := 0; i < 4; i++ {
		dir = filepath.Dir(dir)
	}

	test.InitTestResource(dir)
}

func TestDemoCRUD(t *testing.T) {
	ctx := test.Context()
	entity := &demo.Entity{
		BaseEntity: model.BaseEntity{},
		Name:       "test",
		Status:     0,
	}

	d := dao()
	r := d.Insert(ctx, entity)
	t.Log("insert result", r)

	id := r.LastInsertID

	err := d.SelectByID(ctx, id, entity)
	t.Log("SelectByID", err, entity, *entity.ID, *entity.AddTime, *entity.EditTime)

	r = d.UpdateByIDs(ctx, []*mysql.SqlUpdateColumn{
		{
			Name:  "status",
			Value: 1,
		},
	}, id)
	t.Log("update result", r)

	params := &mysql.SqlQueryParams{
		CondItems: []*mysql.SqlColQueryItem{
			{
				Name:      "id",
				Condition: mysql.SqlCondEqual,
				Value:     id,
			},
		},
		OrderBy: "",
		Offset:  0,
		Cnt:     0,
	}

	cnt, err := d.SimpleTotalAnd(ctx, params.CondItems...)
	t.Log("total", cnt, err)

	var data []*demo.Entity
	err = d.SimpleQueryAnd(ctx, params, &data)
	t.Log("SimpleQueryAnd", err)
	for i, item := range data {
		t.Log(i, item, *item.ID, *item.AddTime, *item.EditTime)
	}

	err = d.Begin(ctx)
	t.Log("begin", err)
	r = d.DeleteByIDs(ctx, id)
	t.Log("delete result", r)
	err = d.Rollback(ctx)
	t.Log("rollback", err)

	_ = d.Begin(ctx)
	_ = d.DeleteByIDs(ctx, id)
	err = d.Commit(ctx)
	t.Log("commit", err)
}

func dao() demo.Dao {
	return demo.NewDao(test.MysqlClient())
}

func TestTrans(t *testing.T) {
	ctx := test.Context()
	demoDao := factory.DefaultDaoFactory.DemoDao(ctx)
	idGenDao := factory.DefaultDaoFactory.IDGenDao(ctx)

	_ = demoDao.Begin(ctx)
	r := demoDao.Insert(ctx, &demo.Entity{Name: "a"})
	if r.Err != nil {
		_ = demoDao.Rollback(ctx)
		return
	}
	r = idGenDao.Insert(ctx, &idgen.Entity{Name: "demo"})
	if r.Err != nil {
		_ = demoDao.Rollback(ctx)
		return
	}

	_ = demoDao.Commit(ctx)
}
