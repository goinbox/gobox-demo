package redis

import (
	"github.com/goinbox/golog"

	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client := getTestClient()

	reply := client.Do("set", "c", "1")
	t.Log(reply.String())
	reply = client.Do("get", "c")
	t.Log(reply.Int())

	reply = client.DoWithoutLog("set", "d", "1")
	t.Log(reply.String())
	reply = client.DoWithoutLog("get", "d")
	t.Log(reply.Int())

	client.Send("set", "a", "a")
	client.Send("set", "b", "b")
	client.Send("get", "a")
	client.Send("get", "b")
	replies, errIndexes := client.ExecPipelining()
	t.Log(errIndexes)
	for _, reply := range replies {
		t.Log(reply.String())
		t.Log(reply.Err)
	}

	client.BeginTrans()
	client.Send("set", "a", "1")
	client.Send("set", "b", "2")
	client.Send("get", "a")
	client.Send("get", "b")
	replies, _ = client.ExecTrans()
	for _, reply := range replies {
		t.Log(reply.String())
		t.Log(reply.Err)
	}

	client.Free()
}

func TestAutoReconnect(t *testing.T) {
	client := getTestClient()

	reply := client.Do("set", "a", "1")
	t.Log(reply.String())
	time.Sleep(time.Second * 4) //set redis-server timeout = 3
	reply = client.Do("get", "a")
	t.Log(reply.Err)
	t.Log(reply.Int())

	time.Sleep(time.Second * 4)

	client.Send("set", "a", "a")
	client.Send("set", "b", "b")
	client.Send("get", "a")
	client.Send("get", "b")
	replies, errIndexes := client.ExecPipelining()
	t.Log(errIndexes)
	for _, reply := range replies {
		t.Log(reply.String())
		t.Log(reply.Err)
	}

	time.Sleep(time.Second * 4)

	client.BeginTrans()
	client.Send("set", "a", "1")
	client.Send("set", "b", "2")
	client.Send("get", "a")
	client.Send("get", "b")
	replies, _ = client.ExecTrans()
	for _, reply := range replies {
		t.Log(reply.String())
		t.Log(reply.Err)
	}

	client.Free()
}

func getTestClient() *Client {
	logger, _ := golog.NewSimpleLogger(golog.NewStdoutWriter(), golog.LEVEL_DEBUG, golog.NewConsoleFormater())
	config := NewConfig("127.0.0.1", "6379", "123")

	return NewClient(config, logger)
}
