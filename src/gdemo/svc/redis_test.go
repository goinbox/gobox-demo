package svc

import (
	"gdemo/resource"

	"github.com/goinbox/gomisc"

	"testing"
	"time"
)

func TestRedisSvcSetGetJson(t *testing.T) {
	rs := NewRedisSvc([]byte("TestRedisSvcSetGetJson"), resource.RedisClientPoolList[0])

	item := &demoEntity{
		SqlBaseEntity: SqlBaseEntity{
			Id:       1,
			AddTime:  time.Now().Format(gomisc.TimeGeneralLayout()),
			EditTime: time.Now().Format(gomisc.TimeGeneralLayout()),
		},
		Name:   "tdj",
		Status: 1,
	}

	key := "test_demo_json"
	err := rs.SaveJsonData(key, item, 10)
	if err != nil {
		t.Error(err)
	}

	item = &demoEntity{}
	find, err := rs.GetJsonData(key, item)
	if !find {
		t.Error("not found")
	}
	if err != nil {
		t.Error(err)
	}

	t.Log(item)
}

func TestRedisSvcSetGetHash(t *testing.T) {
	rs := NewRedisSvc([]byte("TestRedisSvcSetGetHash"), resource.RedisClientPoolList[0])

	item := &demoEntity{
		SqlBaseEntity: SqlBaseEntity{
			Id:       1,
			AddTime:  time.Now().Format(gomisc.TimeGeneralLayout()),
			EditTime: time.Now().Format(gomisc.TimeGeneralLayout()),
		},
		Name:   "tdj",
		Status: 1,
	}

	key := "test_demo_hash"
	err := rs.SaveHashEntity(key, item, 10)
	if err != nil {
		t.Error(err)
	}

	item = &demoEntity{}
	find, err := rs.GetHashEntity(key, item)
	if !find {
		t.Error("not found")
	}
	if err != nil {
		t.Error(err)
	}

	t.Log(item)
}
