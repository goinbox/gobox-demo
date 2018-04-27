package redis

import (
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	config := &PConfig{NewClientFunc: newRedisTestClient}
	config.Size = 100
	config.MaxIdleTime = time.Second * 5
	config.KeepAliveInterval = time.Second * 3

	pool := NewPool(config)
	testPool(pool, t)

	config.LogKeepAlive = true
	pool = NewPool(config)
	testPool(pool, t)
}

func newRedisTestClient() (*Client, error) {
	return getTestClient(), nil
}

func testPool(pool *Pool, t *testing.T) {
	client, _ := pool.Get()
	client.Do("set", "redis_pool", "pool_test")
	reply := client.Do("get", "redis_pool")
	t.Log(reply.String())
	pool.Put(client)

	time.Sleep(time.Second * 4)
	client, _ = pool.Get()
	reply = client.Do("get", "redis_pool")
	t.Log(reply.String())
	pool.Put(client)
}
