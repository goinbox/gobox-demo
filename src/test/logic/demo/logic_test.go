package demo

import (
	"testing"

	"github.com/goinbox/mysql"

	demoLogic "gdemo/logic/demo"
	"gdemo/logic/factory"
	"gdemo/model"
	"gdemo/model/demo"
	"gdemo/test"
)

func init() {
	test.InitMysql()
}

func TestDemoCRUD(t *testing.T) {
	name := "test"
	entity := &demo.Entity{
		BaseEntity: model.BaseEntity{},
		Name:       name,
		Status:     0,
	}

	l := logic()
	ctx := test.Context()

	err := l.Insert(ctx, entity)
	t.Log("Insert", err)

	conds := []*mysql.SqlColQueryItem{
		{
			Name:      "name",
			Condition: mysql.SqlCondEqual,
			Value:     name,
		},
	}
	params := &mysql.SqlQueryParams{
		CondItems: conds,
		OrderBy:   "",
		Offset:    0,
		Cnt:       0,
	}

	total, err := l.SimpleTotalAnd(ctx, conds...)
	t.Log("SimpleTotalAnd", total, err)

	data, err := l.SimpleQueryAnd(ctx, params)
	item := data[0]
	t.Log("SimpleQueryAnd", item, err)

	err = l.UpdateByIDs(ctx, map[string]interface{}{"status": 1}, *item.ID)
	t.Log("UpdateByIDs", err)

	entity, err = l.SelectByID(ctx, *item.ID)
	t.Log("SelectByID", entity, err)

	err = l.DeleteByIDs(ctx, *entity.ID)
}

func logic() demoLogic.Logic {
	return factory.DefaultLogicFactory.DemoLogic()
}
