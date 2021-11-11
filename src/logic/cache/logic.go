package cache

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/goinbox/redis"

	"gdemo/pcontext"
	"gdemo/resource"
)

type Logic interface {
	Set(key string, value interface{}, expireSeconds int64) error
	Get(key string, value interface{}) (bool, error)
	Del(key string) error

	SetStruct(key string, value interface{}, expireSeconds int64) error
	GetStruct(key string, value interface{}) (bool, error)
}

type logic struct {
	ctx *pcontext.Context
}

func NewLogic(ctx *pcontext.Context) *logic {
	return &logic{ctx: ctx}
}

func (l *logic) Set(key string, value interface{}, expireSeconds int64) error {
	client := l.client()

	var err error
	if expireSeconds > 0 {
		err = client.Do("set", key, value, "ex", expireSeconds).Err
	} else {
		err = client.Do("set", key, value).Err
	}

	if err != nil {
		return fmt.Errorf("CacheLogic.Set client.Do error: %w", err)
	}

	return nil
}

func (l *logic) Get(key string, value interface{}) (bool, error) {
	reply := l.get(key)
	if reply.Err != nil {
		if reply.Nil() {
			return false, nil
		}
		return false, fmt.Errorf("CacheLogic.Get get error: %w", reply.Err)
	}

	v := reflect.ValueOf(value).Elem()
	v.Set(reflect.ValueOf(reply.Value()))

	return true, nil
}

func (l *logic) get(key string) *redis.Reply {
	return l.client().Do("get", key)
}

func (l *logic) Del(key string) error {
	client := l.client()

	reply := client.Do("del", key)
	if reply.Err != nil {
		return fmt.Errorf("CacheLogic.Del client.Do error: %w", reply.Err)
	}

	return nil
}

func (l *logic) SetStruct(key string, value interface{}, expireSeconds int64) error {
	content, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("CacheLogic.SetStruct json.Marshal error: %w", err)
	}

	return l.Set(key, content, expireSeconds)
}

func (l *logic) GetStruct(key string, value interface{}) (bool, error) {
	reply := l.get(key)
	if reply.Err != nil {
		if reply.Nil() {
			return false, nil
		}
		return false, fmt.Errorf("CacheLogic.GetStruct get error: %w", reply.Err)
	}

	content, err := reply.String()
	if err != nil {
		return false, fmt.Errorf("CacheLogic.GetStruct reply.String error: %w", err)
	}

	err = json.Unmarshal([]byte(content), value)
	if err != nil {
		return false, fmt.Errorf("CacheLogic.GetStruct json.Unmarshal error: %w", err)
	}

	return true, nil
}

func (l *logic) client() *redis.Client {
	return resource.RedisClient(l.ctx.Logger)
}
