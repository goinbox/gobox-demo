package store

import (
	"testing"

	"gdemo/define/entity"
	"gdemo/resource"
)

func TestSqlRedisBindStore(t *testing.T) {
	srs := &SqlRedisBindStore{
		SqlStore:   NewSqlStore([]byte("TestSqlRedisBindSvc"), resource.MysqlClientPool, true),
		RedisStore: NewRedisStore([]byte("TestSqlRedisBindSvc"), resource.RedisClientPoolList[0]),
	}

	tableName, entityName := "demo", "demo"
	redisKeyPrefix := "test_sql_redis_bind"

	var merr, rerr error
	var ids []int64
	var find bool

	resource.TestLogger.Notice([]byte("test Insert"))

	item := &entity.DemoEntity{
		Name:   "tdj",
		Status: 1,
	}

	ids, merr, rerr = srs.Insert(tableName, entityName, redisKeyPrefix, entity.DemoColNames, 10, item)
	t.Log(ids, merr, rerr)

	resource.TestLogger.Notice([]byte("test GetById"))

	item = &entity.DemoEntity{}
	find, merr, rerr = srs.GetById(tableName, entityName, redisKeyPrefix, ids[0], 10, item)
	t.Log(find, merr, rerr, item)

	resource.TestLogger.Notice([]byte("test UpdateById"))

	newDemo := &entity.DemoEntity{
		Name: "new-demo",
	}
	updateFields := map[string]bool{"name": true}
	setItems, merr, rerr := srs.UpdateById(tableName, entityName, redisKeyPrefix, ids[0], newDemo, updateFields, 10)
	t.Log(merr, rerr)
	for i, item := range setItems {
		t.Log(i, item)
	}
	item = &entity.DemoEntity{}
	find, merr, rerr = srs.GetById(tableName, entityName, redisKeyPrefix, ids[0], 10, item)
	t.Log(find, merr, rerr, item)

	resource.TestLogger.Notice([]byte("test TotalRows"))

	total, merr, rerr := srs.TotalRows(tableName, redisKeyPrefix, 10)
	t.Log(total, merr, rerr)

	resource.TestLogger.Notice([]byte("test TotalRows"))

	find, merr, rerr = srs.DeleteById(tableName, entityName, redisKeyPrefix, ids[0])
	t.Log(find, merr, rerr)
}
