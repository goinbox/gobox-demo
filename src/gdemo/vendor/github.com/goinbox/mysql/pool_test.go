package mysql

import (
	"database/sql"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	config := &PConfig{NewClientFunc: newMysqlTestClient}
	config.Size = 100
	config.MaxIdleTime = time.Second * 5

	pool := NewPool(config)

	testPool(pool, t)
	testPool(pool, t)

	time.Sleep(time.Second * 7)
	testPool(pool, t)
}

func newMysqlTestClient() (*Client, error) {
	return NewClient(getTestConfig(), nil)
}

func testPool(pool *Pool, t *testing.T) {
	client, _ := pool.Get()
	row := client.QueryRow("SELECT * FROM demo WHERE id = ?", 1)
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

	pool.Put(client)
}
