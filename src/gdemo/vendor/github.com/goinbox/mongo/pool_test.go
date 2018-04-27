package mongo

import (
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func TestPool(t *testing.T) {
	config := &PConfig{NewClientFunc: newTestClient}
	config.Size = 100
	config.MaxIdleTime = time.Second * 5

	pool := NewPool(config)
	testPool(pool, t)
}

func newTestClient() (*Client, error) {
	return getTestClient(), nil
}

func testPool(pool *Pool, t *testing.T) {
	client, _ := pool.Get()

	query := NewQuery()
	result := []bson.M{}
	err := client.Query(MONGO_TEST_COLLECTION, query).All(&result)
	if err != nil {
		t.Error(err)
	}
	jsonData, _ := json.Marshal(result)
	t.Logf("%s", jsonData)

	pool.Put(client)
}
