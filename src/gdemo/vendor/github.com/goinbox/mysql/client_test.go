package mysql

import (
	"github.com/goinbox/golog"

	"database/sql"
	"strconv"
	"testing"
)

type tableDemoRowItem struct {
	Id       int64
	AddTime  string
	EditTime string
	Name     string
	Status   int
}

var client *Client

func init() {
	config := getTestConfig()
	client = getTestClient(config)
}

func TestClientExec(t *testing.T) {
	result, err := client.Exec("INSERT INTO demo (name) VALUES (?)", "a")
	if err != nil {
		t.Log("exec error: " + err.Error())
	} else {
		li, err := result.LastInsertId()
		if err != nil {
			t.Log("lastInsertId error: " + err.Error())
		} else {
			t.Log("lastInsertId: " + strconv.FormatInt(li, 10))
		}

		rf, err := result.RowsAffected()
		if err != nil {
			t.Log("rowsAffected error: " + err.Error())
		} else {
			t.Log("rowsAffected: " + strconv.FormatInt(rf, 10))
		}
	}
}

func TestClientQuery(t *testing.T) {
	rows, err := client.Query("SELECT * FROM demo WHERE id IN (?,?)", 1, 5)
	if err != nil {
		t.Log("query error: " + err.Error())
	} else {
		for rows.Next() {
			item := new(tableDemoRowItem)
			err = rows.Scan(&item.Id, &item.AddTime, &item.EditTime, &item.Name, &item.Status)
			if err != nil {
				t.Log("rows scan error: " + err.Error())
			} else {
				t.Log(item)
			}
		}
	}
}

func TestClientQueryRow(t *testing.T) {
	row := client.QueryRow("SELECT * FROM demo WHERE id = ?", 5)
	item := new(tableDemoRowItem)
	err := row.Scan(&item.Id, &item.AddTime, &item.EditTime, &item.Name, &item.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			t.Log("no rows: " + err.Error())
		} else {
			t.Log("row scan error: " + err.Error())
		}
	} else {
		t.Log(item)
	}
}

func TestClientTrans(t *testing.T) {
	client.Begin()

	row := client.QueryRow("SELECT * FROM demo WHERE id = ?", 1)
	item := new(tableDemoRowItem)
	err := row.Scan(&item.Id, &item.AddTime, &item.EditTime, &item.Name, &item.Status)
	if err != nil {
		t.Log("row scan error: " + err.Error())
	} else {
		t.Log(item)
	}

	client.Commit()

	err = client.Rollback()
	t.Log(err)
}

func getTestClient(config *Config) *Client {
	w, _ := golog.NewFileWriter("/tmp/test_mysql.log")
	logger, _ := golog.NewSimpleLogger(w, golog.LEVEL_INFO, golog.NewSimpleFormater())

	client, _ := NewClient(config, logger)

	return client
}

func getTestConfig() *Config {
	return NewConfig("root", "123", "127.0.0.1", "3306", "gobox-demo")
}
