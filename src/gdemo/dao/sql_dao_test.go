package dao

import (
	"gdemo/misc"

	"github.com/goinbox/gomisc"

	"testing"
	"time"
)

const (
	SQL_TEST_TABLE_NAME = "demo"
)

type sqlTestEntity struct {
	Id       int64
	AddTime  string
	EditTime string
	Name     string
	Status   int
}

func TestSqlDaoRead(t *testing.T) {
	dao := &SqlDao{misc.MysqlTestClient()}
	entity := new(sqlTestEntity)

	row := dao.SelectById(SQL_TEST_TABLE_NAME, "*", 1)
	row.Scan(&entity.Id, &entity.AddTime, &entity.EditTime, &entity.Name, &entity.Status)
	t.Log(entity)

	rows, _ := dao.SelectByIds(SQL_TEST_TABLE_NAME, "*", "id desc", 1)
	for rows.Next() {
		rows.Scan(&entity.Id, &entity.AddTime, &entity.EditTime, &entity.Name, &entity.Status)
		t.Log("entity", entity)
	}

	condItems := []*SqlColQueryItem{
		NewSqlColQueryItem("name", SQL_COND_LIKE, "%a%"),
		NewSqlColQueryItem("id", SQL_COND_BETWEEN, []int64{0, 100}),
		NewSqlColQueryItem("status", SQL_COND_EQUAL, 0),
	}
	rows, _ = dao.SimpleQueryAnd(SQL_TEST_TABLE_NAME, "*", "id desc", 0, 10, condItems...)
	for rows.Next() {
		rows.Scan(&entity.Id, &entity.AddTime, &entity.EditTime, &entity.Name, &entity.Status)
		t.Log(entity)
	}

	total, _ := dao.SimpleTotalAnd(SQL_TEST_TABLE_NAME, condItems...)
	t.Log(total)
}

func TestSqlDaoWrite(t *testing.T) {
	dao := &SqlDao{misc.MysqlTestClient()}

	var colNames = []string{"id", "add_time", "edit_time", "name", "status"}
	var colsValues [][]interface{}

	ts := time.Now().Format(gomisc.TimeGeneralLayout())
	for i, name := range []string{"a", "b", "c"} {
		colValues := []interface{}{
			int64(i + 10),
			ts,
			ts,
			name,
			i % 10,
		}
		colsValues = append(colsValues, colValues)
	}
	result := dao.Insert(SQL_TEST_TABLE_NAME, colNames, colsValues...)
	t.Log(result)

	id := result.LastInsertId
	setItems := []*SqlColQueryItem{
		NewSqlColQueryItem("name", "", "abc"),
		NewSqlColQueryItem("edit_time", "", ts),
	}
	result = dao.UpdateById(SQL_TEST_TABLE_NAME, id, setItems...)
	t.Log(result)

	result = dao.DeleteById(SQL_TEST_TABLE_NAME, id)
	t.Log(result)
}
