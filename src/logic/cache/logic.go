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
	Set(ctx *pcontext.Context, key string, value interface{}, expireSeconds int64) error
	Get(ctx *pcontext.Context, key string, value interface{}) (bool, error)
	Del(ctx *pcontext.Context, key string) error

	SetStruct(ctx *pcontext.Context, key string, value interface{}, expireSeconds int64) error
	GetStruct(ctx *pcontext.Context, key string, value interface{}) (bool, error)
}

type logic struct {
}

func NewLogic() *logic {
	return &logic{}
}

func (l *logic) Set(ctx *pcontext.Context, key string, value interface{}, expireSeconds int64) error {
	client := l.client(ctx)

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

func (l *logic) Get(ctx *pcontext.Context, key string, value interface{}) (bool, error) {
	reply := l.get(ctx, key)
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

func (l *logic) get(ctx *pcontext.Context, key string) *redis.Reply {
	return l.client(ctx).Do("get", key)
}

func (l *logic) Del(ctx *pcontext.Context, key string) error {
	client := l.client(ctx)

	reply := client.Do("del", key)
	if reply.Err != nil {
		return fmt.Errorf("CacheLogic.Del client.Do error: %w", reply.Err)
	}

	return nil
}

func (l *logic) SetStruct(ctx *pcontext.Context, key string, value interface{}, expireSeconds int64) error {
	content, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("CacheLogic.SetStruct json.Marshal error: %w", err)
	}

	return l.Set(ctx, key, content, expireSeconds)
}

func (l *logic) GetStruct(ctx *pcontext.Context, key string, value interface{}) (bool, error) {
	reply := l.get(ctx, key)
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

func (l *logic) client(ctx *pcontext.Context) *redis.Client {
	return resource.RedisClient(ctx.Logger)
}
