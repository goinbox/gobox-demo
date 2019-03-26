package demo

import (
	"gdemo/conf"
	"gdemo/resource"
	"gdemo/svc"

	"github.com/goinbox/golog"

	"reflect"
)

const (
	DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS     = 60 * 30
	DEF_DEMO_TOTAL_ROWS_CACHE_EXPIRE_SECONDS = 60 * 10
)

var demoEntityType reflect.Type = reflect.TypeOf(DemoEntity{})
var demoColNames []string = svc.ReflectColNames(demoEntityType)

type DemoEntity struct {
	svc.SqlBaseEntity

	Name   string `mysql:"name" json:"name" redis:"name"`
	Status int    `mysql:"status" json:"status" redis:"status"`
}

type DemoSvc struct {
	*svc.BaseSvc
	*svc.SqlRedisBindSvc

	EntityName     string
	RedisKeyPrefix string
}

func NewDemoSvc(traceId []byte, logger golog.ILogger) *DemoSvc {
	return &DemoSvc{
		BaseSvc: &svc.BaseSvc{
			TraceId:      traceId,
			AccessLogger: logger,
		},
		SqlRedisBindSvc: &svc.SqlRedisBindSvc{
			SqlSvc:   svc.NewSqlSvc(logger, resource.MysqlClientPool, true),
			RedisSvc: svc.NewRedisSvc(logger, resource.RedisClientPoolList[0]),
		},

		EntityName:     "demo",
		RedisKeyPrefix: conf.BaseConf.PrjName,
	}
}

func (d *DemoSvc) Insert(entities ...*DemoEntity) ([]int64, error) {
	is := make([]interface{}, len(entities))
	for i, entity := range entities {
		is[i] = entity
	}

	ids, merr, rerr := d.SqlRedisBindSvc.Insert(d.EntityName, d.EntityName, d.RedisKeyPrefix, demoColNames, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS, is...)
	if rerr != nil {
		d.ErrorLog([]byte("DemoSvc.Insert"), []byte(rerr.Error()))
	}

	return ids, merr
}

func (d *DemoSvc) GetById(id int64) (*DemoEntity, error) {
	entity := new(DemoEntity)

	find, merr, rerr := d.SqlRedisBindSvc.GetById(d.EntityName, d.EntityName, d.RedisKeyPrefix, id, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS, entity)
	if rerr != nil {
		d.ErrorLog([]byte("DemoSvc.GetById"), []byte(rerr.Error()))
	}

	if merr != nil {
		return nil, merr
	}
	if !find {
		return nil, nil
	}

	return entity, nil
}

func (d *DemoSvc) DeleteById(id int64) (bool, error) {
	find, merr, rerr := d.SqlRedisBindSvc.DeleteById(d.EntityName, d.EntityName, d.RedisKeyPrefix, id)
	if rerr != nil {
		d.ErrorLog([]byte("DemoSvc.DeleteById"), []byte(rerr.Error()))
	}

	return find, merr
}

func (d *DemoSvc) UpdateById(id int64, newEntity *DemoEntity, updateFields map[string]bool) (bool, error) {
	setItems, merr, rerr := d.SqlRedisBindSvc.UpdateById(d.EntityName, d.EntityName, d.RedisKeyPrefix, id, newEntity, updateFields, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS)
	if rerr != nil {
		d.ErrorLog([]byte("DemoSvc.UpdateById"), []byte(rerr.Error()))
	}

	return (setItems != nil), merr
}

func (d *DemoSvc) ListByIds(ids ...int64) ([]*DemoEntity, error) {
	var entities []*DemoEntity

	err := d.SqlSvc.ListByIds(d.EntityName, ids, "id desc", demoEntityType, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (d *DemoSvc) SimpleQueryAnd(sqp *svc.SqlQueryParams) ([]*DemoEntity, error) {
	var entities []*DemoEntity

	err := d.SqlSvc.SimpleQueryAnd(d.EntityName, sqp, demoEntityType, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}
