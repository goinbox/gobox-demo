package store

import (
	"testing"
	"time"

	"github.com/goinbox/gomisc"

	"gdemo/define/entity"
	"gdemo/resource"
)

func TestRedisStoreSetGetJson(t *testing.T) {
	rs := NewRedisStore([]byte("TestRedisStoreSetGetJson"), resource.RedisClientPoolList[0])

	item := &entity.DemoEntity{
		SqlBaseEntity: entity.SqlBaseEntity{
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

	item = &entity.DemoEntity{}
	find, err := rs.GetJsonData(key, item)
	if !find {
		t.Error("not found")
	}
	if err != nil {
		t.Error(err)
	}

	t.Log(item)
}

func TestRedisStoreSetGetHash(t *testing.T) {
	rs := NewRedisStore([]byte("TestRedisStoreSetGetHash"), resource.RedisClientPoolList[0])

	item := &entity.DemoEntity{
		SqlBaseEntity: entity.SqlBaseEntity{
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

	item = &entity.DemoEntity{}
	find, err := rs.GetHashEntity(key, item)
	if !find {
		t.Error("not found")
	}
	if err != nil {
		t.Error(err)
	}

	t.Log(item)
}
