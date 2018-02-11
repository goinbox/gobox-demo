package svc

import (
	"github.com/goinbox/redis"

	"encoding/json"
	"errors"
	"reflect"
)

const (
	ENTITY_REDIS_HASH_FIELD_TAG = "redis"
)

type RedisBaseSvc struct {
	*BaseSvc

	rclient *redis.Client
}

func NewRedisBaseSvc(bs *BaseSvc, rclient *redis.Client) *RedisBaseSvc {
	return &RedisBaseSvc{
		BaseSvc: bs,

		rclient: rclient,
	}
}

func (r *RedisBaseSvc) saveJsonDataToRedis(key string, v interface{}, expireSeconds int64) error {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		r.elogger.Warning([]byte("json_encode " + key + " error: " + err.Error()))
		return err
	}

	err = r.rclient.Do("set", key, string(jsonBytes), "ex", expireSeconds).Err
	if err != nil {
		r.rclient.Free()
		r.elogger.Warning([]byte("set " + key + " to redis error: " + err.Error()))
		return err
	}

	return nil
}

func (r *RedisBaseSvc) getJsonDataFromRedis(key string, v interface{}) (bool, error) {
	reply := r.rclient.Do("get", key)
	if reply.Err != nil {
		r.rclient.Free()
		r.elogger.Warning([]byte("get " + key + " from redis error: " + reply.Err.Error()))
		return false, reply.Err
	}

	if reply == nil {
		return false, nil
	}

	jsonBytes, err := reply.Bytes()
	if err != nil {
		r.elogger.Warning([]byte("reply " + key + " from redis error: " + err.Error()))
		return false, err
	}

	err = json.Unmarshal(jsonBytes, v)
	if err != nil {
		r.elogger.Warning([]byte("json_decode " + key + " from redis error: " + err.Error()))
		return false, err
	}

	return true, nil
}

func (r *RedisBaseSvc) saveHashEntityToRedis(key string, entityPtr interface{}, expireSeconds int64) error {
	eargs := r.reflectSaveHashEntityArgs(reflect.ValueOf(entityPtr).Elem())
	args := make([]interface{}, len(eargs)+1)
	args[0] = key
	for i, arg := range eargs {
		args[i+1] = arg
	}

	r.rclient.Send("hmset", args...)
	if expireSeconds > 0 {
		r.rclient.Send("expire", key, expireSeconds)
	}
	replies, errIndexes := r.rclient.ExecPipelining()
	if len(errIndexes) != 0 {
		r.rclient.Free()
		msg := "hmset " + key + " to redis error:"
		for _, i := range errIndexes {
			msg += " " + replies[i].Err.Error()
		}
		r.elogger.Warning([]byte(msg))
		return errors.New(msg)
	}

	return nil
}

func (r *RedisBaseSvc) getHashEntityFromRedis(key string, entityPtr interface{}) (bool, error) {
	reply := r.rclient.Do("hgetall", key)
	if reply.Err != nil {
		r.rclient.Free()
		r.elogger.Warning([]byte("hgetall " + key + " from redis error: " + reply.Err.Error()))
		return false, reply.Err
	}

	if reply.ArrReplyIsNil() {
		return false, nil
	}

	err := reply.Struct(entityPtr)
	if err != nil {
		r.elogger.Warning([]byte("reply to struct " + key + " from redis error: " + err.Error()))
		return false, err
	}

	return true, nil
}

func (r *RedisBaseSvc) reflectSaveHashEntityArgs(rev reflect.Value) []interface{} {
	var args []interface{}
	ret := rev.Type()

	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct {
			args = r.reflectSaveHashEntityArgs(revf)
			continue
		}

		retf := ret.Field(i)
		fn, ok := retf.Tag.Lookup(ENTITY_REDIS_HASH_FIELD_TAG)
		if ok {
			args = append(args, fn, revf.Interface())
		}
	}

	return args
}
