package cache

import (
	"testing"
	"time"

	"gdemo/logic/cache"
	"gdemo/test"
)

func init() {
	test.InitRedis()
}

func TestCacheLogicGetSet(t *testing.T) {
	ctx := test.Context()
	logic := cache.NewLogic()

	key := "test"
	err := logic.Set(ctx, key, "cache logic", 2)
	t.Log("set", err)

	var v string
	exist, err := logic.Get(ctx, key, &v)
	t.Log("get before expire", exist, err, v)

	time.Sleep(time.Second * 3)
	exist, err = logic.Get(ctx, key, &v)
	t.Log("get after expire", exist, err)

	_ = logic.Set(ctx, key, "cache logic", 3)
	err = logic.Del(ctx, key)
	t.Log("del", err)

	exist, err = logic.Get(ctx, key, &v)
	t.Log("get after del", exist, err)
}

func TestCacheLogicGetSetStruct(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}

	ctx := test.Context()
	logic := cache.NewLogic()

	key := "test"
	err := logic.SetStruct(ctx, key, &person{
		Name: "zhangsan",
		Age:  10,
	}, 0)
	t.Log("set struct")

	item := new(person)

	exist, err := logic.GetStruct(ctx, key, item)
	t.Log("get struct", exist, err, item)
}
