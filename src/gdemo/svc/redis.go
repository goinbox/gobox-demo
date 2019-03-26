package svc

import (
	"gdemo/resource"
	"github.com/goinbox/golog"
	"github.com/goinbox/redis"

	"encoding/json"
	"errors"
	"reflect"
)

const (
	ENTITY_REDIS_HASH_FIELD_TAG = "redis"
)

type RedisSvc struct {
	logger golog.ILogger
	pool   *redis.Pool

	client *redis.Client
}

func NewRedisSvc(logger golog.ILogger, pool *redis.Pool) *RedisSvc {
	return &RedisSvc{
		logger: logger,
		pool:   pool,
	}
}

func (r *RedisSvc) Client() *redis.Client {
	if r.client == nil {
		r.client, _ = r.pool.Get()
		r.client.SetLogger(r.logger)
	}

	return r.client
}

func (r *RedisSvc) SendBackClient() {
	if r.client.Connected() {
		r.client.SetLogger(resource.NoopLogger)
		_ = r.pool.Put(r.client)
	}

	r.client = nil
}

func (r *RedisSvc) Renew(logger golog.ILogger, pool *redis.Pool) *RedisSvc {
	if r.client != nil {
		r.SendBackClient()
	}

	r.logger = logger
	r.pool = pool

	return r
}

func (r *RedisSvc) SetPool(pool *redis.Pool) *RedisSvc {
	return r.Renew(r.logger, pool)
}

func (r *RedisSvc) SaveJsonData(key string, v interface{}, expireSeconds int64) error {
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

func (r *RedisSvc) GetJsonData(key string, v interface{}) (bool, error) {
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

func (r *RedisSvc) SaveHashEntity(key string, entityPtr interface{}, expireSeconds int64) error {
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

func (r *RedisSvc) GetHashEntity(key string, entityPtr interface{}) (bool, error) {
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

func (r *RedisSvc) ReflectSaveHashEntityArgs(rev reflect.Value) []interface{} {
	var args []interface{}
	ret := rev.Type()

	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct {
			args = r.ReflectSaveHashEntityArgs(revf)
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
