package svc

import (
	"errors"
	"gdemo/dao"
	"strconv"
)

type SqlRedisBindSvc struct {
	*BaseSvc
	*SqlBaseSvc
	*RedisBaseSvc

	redisKeyPrefix string
}

func NewSqlRedisBindSvc(bs *BaseSvc, ss *SqlBaseSvc, rs *RedisBaseSvc, redisKeyPrefix string) *SqlRedisBindSvc {
	return &SqlRedisBindSvc{bs, ss, rs, redisKeyPrefix}
}

func (s *SqlRedisBindSvc) redisKeyForEntity(id int64) string {
	return s.redisKeyPrefix + "_entity_" + s.entityName + "_id_" + strconv.FormatInt(id, 10)
}

func (s *SqlRedisBindSvc) Insert(tableName string, colNames []string, expireSeconds int64, entities ...interface{}) ([]int64, error) {
	ids, err := s.SqlBaseSvc.Insert(tableName, colNames, entities...)
	if err != nil {
		return nil, err
	}

	for i, id := range ids {
		s.saveHashEntityToRedis(s.redisKeyForEntity(id), entities[i], expireSeconds)
	}

	return ids, nil
}

func (s *SqlRedisBindSvc) GetById(tableName string, id, expireSeconds int64, entityPtr interface{}) (bool, error) {
	rk := s.redisKeyForEntity(id)

	find, err := s.getHashEntityFromRedis(rk, entityPtr)
	if err != nil {
		return s.SqlBaseSvc.GetById(tableName, id, entityPtr)
	}
	if find {
		return true, nil
	}

	find, err = s.SqlBaseSvc.GetById(tableName, id, entityPtr)
	if err != nil {
		s.elogger.Error([]byte("getById from mysql error"))
		return false, err
	}
	if !find {
		return false, nil
	}

	s.saveHashEntityToRedis(rk, entityPtr, expireSeconds)

	return true, nil
}

func (s *SqlRedisBindSvc) DeleteById(tableName string, id int64) (bool, error) {
	result := s.dao.DeleteById(tableName, id)
	if result.Err != nil {
		s.elogger.Error([]byte("delete from mysql error: " + result.Err.Error()))
		return false, result.Err
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	rk := s.redisKeyForEntity(id)
	err := s.rclient.Do("del", rk).Err
	if err != nil {
		s.elogger.Warning([]byte("del key " + rk + " from redis failed: " + err.Error()))
	}

	return true, nil
}

func (s *SqlRedisBindSvc) UpdateById(tableName string, id int64, newEntityPtr interface{}, updateFields map[string]bool, expireSeconds int64) (bool, error) {
	setItems, err := s.SqlBaseSvc.UpdateById(tableName, id, newEntityPtr, updateFields)

	if err != nil {
		return false, err
	}
	if setItems == nil {
		return false, nil
	}

	s.updateSqlHashEntity(s.redisKeyForEntity(id), setItems, expireSeconds)

	return true, nil
}

func (s *SqlRedisBindSvc) updateSqlHashEntity(key string, setItems []*dao.SqlColQueryItem, expireSeconds int64) error {
	cnt := len(setItems)*2 + 1
	args := make([]interface{}, cnt)
	args[0] = key

	for si, ai := 0, 1; ai < cnt; si++ {
		args[ai] = setItems[si].Name
		ai++
		args[ai] = setItems[si].Value
		ai++
	}

	s.rclient.Send("hmset", args...)
	if expireSeconds > 0 {
		s.rclient.Send("expire", expireSeconds)
	}
	replies, errIndexes := s.rclient.ExecPipelining()
	if len(errIndexes) != 0 {
		s.rclient.Free()
		msg := "hmset key " + key + " to redis error:"
		for _, i := range errIndexes {
			msg += " " + replies[i].Err.Error()
		}
		s.elogger.Warning([]byte(msg))
		return errors.New(msg)
	}

	return nil
}
