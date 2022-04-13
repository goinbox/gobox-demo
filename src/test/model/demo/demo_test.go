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
	entity := &demo.Entity{
		BaseEntity: model.BaseEntity{},
		Name:       "test",
		Status:     0,
	}

	d := dao()
	r := d.Insert(entity)
	t.Log("insert result", r)

	id := r.LastInsertID

	err := d.SelectByID(id, entity)
	t.Log("SelectByID", err, entity, *entity.ID, *entity.AddTime, *entity.EditTime)

	r = d.UpdateByIDs([]*mysql.SqlUpdateColumn{
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

	cnt, err := d.SimpleTotalAnd(params.CondItems...)
	t.Log("total", cnt, err)

	var data []*demo.Entity
	err = d.SimpleQueryAnd(params, &data)
	t.Log("SimpleQueryAnd", err)
	for i, item := range data {
		t.Log(i, item, *item.ID, *item.AddTime, *item.EditTime)
	}

	err = d.Begin()
	t.Log("begin", err)
	r = d.DeleteByIDs(id)
	t.Log("delete result", r)
	err = d.Rollback()
	t.Log("rollback", err)

	_ = d.Begin()
	_ = d.DeleteByIDs(id)
	err = d.Commit()
	t.Log("commit", err)
}

func dao() demo.Dao {
	return demo.NewDao(test.MysqlClient())
}

func TestTrans(t *testing.T) {
	ctx := test.Context()
	demoDao := factory.DefaultDaoFactory.DemoDao(ctx)
	idGenDao := factory.DefaultDaoFactory.IDGenDao(ctx)

	_ = demoDao.Begin()
	r := demoDao.Insert(&demo.Entity{Name: "a"})
	if r.Err != nil {
		_ = demoDao.Rollback()
		return
	}
	r = idGenDao.Insert(&idgen.Entity{Name: "demo"})
	if r.Err != nil {
		_ = demoDao.Rollback()
		return
	}

	_ = demoDao.Commit()
}
