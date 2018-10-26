package mysql

import (
	"github.com/goinbox/golog"

	"database/sql"
	"strconv"
	"testing"
)

var client *Client

type tableDemoRowItem struct {
	Id       int64
	AddTime  string
	EditTime string
	Name     string
	Status   int
}

func init() {
	logger, _ := golog.NewSimpleLogger(golog.NewStdoutWriter(), golog.LEVEL_DEBUG, golog.NewConsoleFormater())

	config := NewConfig("root", "123", "127.0.0.1", "3306", "gobox-demo")
	client, _ = NewClient(config, logger)

	client.Exec("DELETE FROM demo")
}

func TestClientExec(t *testing.T) {
	result, err := client.Exec("INSERT INTO demo (name) VALUES (?),(?)", "a", "b")
	if err != nil {
		t.Error("exec error: " + err.Error())
	} else {
		li, err := result.LastInsertId()
		if err != nil {
			t.Error("lastInsertId error: " + err.Error())
		} else {
			t.Log("lastInsertId: " + strconv.FormatInt(li, 10))
		}

		rf, err := result.RowsAffected()
		if err != nil {
			t.Error("rowsAffected error: " + err.Error())
		} else {
			t.Log("rowsAffected: " + strconv.FormatInt(rf, 10))
		}
	}
}

func TestClientQuery(t *testing.T) {
	rows, err := client.Query("SELECT * FROM demo WHERE name IN (?,?)", "a", "b")
	if err != nil {
		t.Error("query error: " + err.Error())
	} else {
		for rows.Next() {
			item := new(tableDemoRowItem)
			err = rows.Scan(&item.Id, &item.AddTime, &item.EditTime, &item.Name, &item.Status)
			if err != nil {
				t.Error("rows scan error: " + err.Error())
			} else {
				t.Log(item)
			}
		}
	}
}

func TestClientQueryRow(t *testing.T) {
	row := client.QueryRow("SELECT * FROM demo WHERE name = ?", "a")
	item := new(tableDemoRowItem)
	err := row.Scan(&item.Id, &item.AddTime, &item.EditTime, &item.Name, &item.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			t.Log("no rows: " + err.Error())
		} else {
			t.Error("row scan error: " + err.Error())
		}
	} else {
		t.Log(item)
	}
}

func TestClientTrans(t *testing.T) {
	client.Begin()

	row := client.QueryRow("SELECT * FROM demo WHERE name = ?", "a")
	item := new(tableDemoRowItem)
	err := row.Scan(&item.Id, &item.AddTime, &item.EditTime, &item.Name, &item.Status)
	if err != nil {
		t.Error("row scan error: " + err.Error())
	} else {
		t.Log(item)
	}

	client.Commit()

	err = client.Rollback()
	t.Log(err)
}
