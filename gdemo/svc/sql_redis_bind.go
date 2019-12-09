package svc

import (
	"github.com/goinbox/mysql"

	"errors"
	"strconv"
)

type SqlRedisBindSvc struct {
	*SqlSvc
	*RedisSvc
}

func (s *SqlRedisBindSvc) RedisKeyForEntity(id int64, prefix, entityName string) string {
	return prefix + "_entity_" + entityName + "_id_" + strconv.FormatInt(id, 10)
}

func (s *SqlRedisBindSvc) RedisKeyForTotalRows(tableName, redisKeyPrefix string) string {
	return redisKeyPrefix + "_total_rows_" + tableName
}

func (s *SqlRedisBindSvc) Insert(tableName, entityName, redisKeyPrefix string, colNames []string, expireSeconds int64, entities ...interface{}) ([]int64, error, error) {
	ids, err := s.SqlSvc.Insert(tableName, entityName, colNames, entities...)
	if err != nil {
		return nil, err, nil
	}

	var rerr error
	for i, id := range ids {
		rerr = s.RedisSvc.SaveHashEntity(s.RedisKeyForEntity(id, redisKeyPrefix, entityName), entities[i], expireSeconds)
		if rerr != nil {
			break
		}
	}

	return ids, nil, rerr
}

func (s *SqlRedisBindSvc) GetById(tableName, entityName, redisKeyPrefix string, id, expireSeconds int64, entityPtr interface{}) (bool, error, error) {
	rk := s.RedisKeyForEntity(id, redisKeyPrefix, entityName)

	find, rerr := s.RedisSvc.GetHashEntity(rk, entityPtr)
	if rerr != nil {
		find, merr := s.SqlSvc.GetById(tableName, id, entityPtr)
		return find, merr, rerr
	}
	if find {
		return true, nil, nil
	}

	find, merr := s.SqlSvc.GetById(tableName, id, entityPtr)
	if merr != nil {
		return false, merr, nil
	}
	if !find {
		return false, nil, nil
	}

	rerr = s.RedisSvc.SaveHashEntity(rk, entityPtr, expireSeconds)

	return true, nil, rerr
}

func (s *SqlRedisBindSvc) DeleteById(tableName, entityName, redisKeyPrefix string, id int64) (bool, error, error) {
	result := s.Dao().DeleteById(tableName, id)
	defer s.SqlSvc.SendBackClient()

	if result.Err != nil {
		return false, result.Err, nil
	}

	if result.RowsAffected == 0 {
		return false, nil, nil
	}

	rk := s.RedisKeyForEntity(id, redisKeyPrefix, entityName)
	err := s.RedisSvc.Client().Do("del", rk).Err
	defer s.RedisSvc.SendBackClient()

	return true, nil, err
}

func (s *SqlRedisBindSvc) UpdateById(tableName, entityName, redisKeyPrefix string, id int64, newEntityPtr interface{}, updateFields map[string]bool, expireSeconds int64) ([]*mysql.SqlColQueryItem, error, error) {
	setItems, err := s.SqlSvc.UpdateById(tableName, id, newEntityPtr, updateFields)

	if err != nil {
		return nil, err, nil
	}
	if setItems == nil {
		return nil, nil, nil
	}

	err = s.RedisSvc.Del(s.RedisKeyForEntity(id, redisKeyPrefix, entityName))

	return setItems, nil, err
}

func (s *SqlRedisBindSvc) TotalRows(tableName, redisKeyPrefix string, expireSeconds int64) (int64, error, error) {
	rk := s.RedisKeyForTotalRows(tableName, redisKeyPrefix)

	rclient := s.RedisSvc.Client()
	defer s.RedisSvc.SendBackClient()

	reply := rclient.Do("get", rk)
	rerr := reply.Err
	if rerr == nil {
		if !reply.SimpleReplyIsNil() {
			total, err := reply.Int64()
			if err == nil {
				return total, nil, nil
			}
			rclient.Do("del", rk)
		}
	}

	total, merr := s.Dao().SimpleTotalAnd(tableName)
	defer s.SqlSvc.SendBackClient()

	if merr != nil {
		return 0, merr, rerr
	}

	reply = rclient.Do("set", rk, total, "ex", expireSeconds)

	return total, nil, rerr
}

func (s *SqlRedisBindSvc) UpdateSqlHashEntity(redisKey string, setItems []*mysql.SqlColQueryItem, expireSeconds int64) error {
	cnt := len(setItems)*2 + 1
	args := make([]interface{}, cnt)
	args[0] = redisKey

	for si, ai := 0, 1; ai < cnt; si++ {
		args[ai] = setItems[si].Name
		ai++
		args[ai] = setItems[si].Value
		ai++
	}

	rclient := s.RedisSvc.Client()
	defer s.RedisSvc.SendBackClient()

	rclient.Send("hmset", args...)
	if expireSeconds > 0 {
		rclient.Send("expire", redisKey, expireSeconds)
	}
	replies, errIndexes := rclient.ExecPipelining()
	if len(errIndexes) != 0 {
		rclient.Free()
		msg := "hmset key " + redisKey + " to redis error:"
		for _, i := range errIndexes {
			msg += " " + replies[i].Err.Error()
		}
		return errors.New(msg)
	}

	return nil
}
