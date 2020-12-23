package store

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/goinbox/redis"

	"gdemo/resource"
)

const (
	EntityRedisHashFieldTag = "redis"
)

type RedisStore struct {
	traceId []byte
	pool    *redis.Pool

	client *redis.Client
}

func NewRedisStore(traceId []byte, pool *redis.Pool) *RedisStore {
	return &RedisStore{
		traceId: traceId,
		pool:    pool,
	}
}

func (r *RedisStore) Client() *redis.Client {
	if r.client == nil {
		r.client, _ = r.pool.Get()
		r.client.SetLogger(resource.AccessLogger).SetTraceId(r.traceId)
	}

	return r.client
}

func (r *RedisStore) SendBackClient() {
	if r.client.Connected() {
		r.client.SetLogger(resource.NoopLogger)
		_ = r.pool.Put(r.client)
	}

	r.client = nil
}

func (r *RedisStore) Renew(traceId []byte, pool *redis.Pool) *RedisStore {
	if r.client != nil {
		r.SendBackClient()
	}

	r.traceId = traceId
	r.pool = pool

	return r
}

func (r *RedisStore) SetPool(pool *redis.Pool) *RedisStore {
	return r.Renew(r.traceId, pool)
}

func (r *RedisStore) SaveJsonData(key string, v interface{}, expireSeconds int64) error {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	client := r.Client()
	defer r.SendBackClient()

	if expireSeconds > 0 {
		err = client.Do("set", key, string(jsonBytes), "ex", expireSeconds).Err
	} else {
		err = client.Do("set", key, string(jsonBytes)).Err
	}

	return err
}

func (r *RedisStore) GetJsonData(key string, v interface{}) (bool, error) {
	reply := r.Client().Do("get", key)
	defer r.SendBackClient()

	if reply.Err != nil {
		return false, reply.Err
	}

	if reply == nil {
		return false, nil
	}

	jsonBytes, err := reply.Bytes()
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(jsonBytes, v)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *RedisStore) SaveHashEntity(key string, entityPtr interface{}, expireSeconds int64) error {
	eargs := r.ReflectSaveHashEntityArgs(reflect.ValueOf(entityPtr).Elem())
	args := make([]interface{}, len(eargs)+1)
	args[0] = key
	for i, arg := range eargs {
		args[i+1] = arg
	}

	client := r.Client()
	defer r.SendBackClient()

	client.Send("hmset", args...)
	if expireSeconds > 0 {
		client.Send("expire", key, expireSeconds)
	}
	replies, errIndexes := client.ExecPipelining()
	if len(errIndexes) != 0 {
		msg := "hmset " + key + " to redis error:"
		for _, i := range errIndexes {
			msg += " " + replies[i].Err.Error()
		}
		return errors.New(msg)
	}

	return nil
}

func (r *RedisStore) GetHashEntity(key string, entityPtr interface{}) (bool, error) {
	reply := r.Client().Do("hgetall", key)
	defer r.SendBackClient()

	if reply.Err != nil {
		return false, reply.Err
	}

	if reply.ArrReplyIsNil() {
		return false, nil
	}

	err := reply.Struct(entityPtr)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *RedisStore) ReflectSaveHashEntityArgs(rev reflect.Value) []interface{} {
	var args []interface{}
	ret := rev.Type()

	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct {
			args = r.ReflectSaveHashEntityArgs(revf)
			continue
		}

		retf := ret.Field(i)
		fn, ok := retf.Tag.Lookup(EntityRedisHashFieldTag)
		if ok {
			args = append(args, fn, revf.Interface())
		}
	}

	return args
}

func (r *RedisStore) Del(key string) error {
	err := r.Client().Do("del", key).Err
	defer r.SendBackClient()

	return err
}
