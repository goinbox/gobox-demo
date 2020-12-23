package demo

import (
	"gdemo/conf"
	"gdemo/define"
	"gdemo/define/entity"
	"gdemo/resource"
	"gdemo/store"
	"gdemo/svc"
)

const (
	DefDemoEntityCacheExpireSeconds    = 60 * 30
	DefDemoTotalRowsCacheExpireSeconds = 60 * 10
)

type DemoSvc struct {
	*svc.BaseSvc
	*store.SqlRedisBindStore

	EntityName     string
	RedisKeyPrefix string
}

func NewDemoSvc(traceId []byte) *DemoSvc {
	return &DemoSvc{
		BaseSvc: &svc.BaseSvc{
			TraceId: traceId,
		},
		SqlRedisBindStore: &store.SqlRedisBindStore{
			SqlStore:   store.NewSqlStore(traceId, resource.MysqlClientPool, true),
			RedisStore: store.NewRedisStore(traceId, resource.RedisClientPoolList[0]),
		},

		EntityName:     "demo",
		RedisKeyPrefix: conf.BaseConf.PrjName,
	}
}

func (d *DemoSvc) Insert(entities ...*entity.DemoEntity) ([]int64, error) {
	is := make([]interface{}, len(entities))
	for i, entity := range entities {
		is[i] = entity
	}

	ids, merr, rerr := d.SqlRedisBindStore.Insert(d.EntityName, d.EntityName, d.RedisKeyPrefix, entity.DemoColNames, DefDemoEntityCacheExpireSeconds, is...)
	if rerr != nil {
		d.ErrorLog([]byte("DemoSvc.Insert"), []byte(rerr.Error()))
	}

	return ids, merr
}

func (d *DemoSvc) GetById(id int64) (*entity.DemoEntity, error) {
	entity := new(entity.DemoEntity)

	find, merr, rerr := d.SqlRedisBindStore.GetById(d.EntityName, d.EntityName, d.RedisKeyPrefix, id, DefDemoEntityCacheExpireSeconds, entity)
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
	find, merr, rerr := d.SqlRedisBindStore.DeleteById(d.EntityName, d.EntityName, d.RedisKeyPrefix, id)
	if rerr != nil {
		d.ErrorLog([]byte("DemoSvc.DeleteById"), []byte(rerr.Error()))
	}

	return find, merr
}

func (d *DemoSvc) UpdateById(id int64, newEntity *entity.DemoEntity, updateFields map[string]bool) (bool, error) {
	setItems, merr, rerr := d.SqlRedisBindStore.UpdateById(d.EntityName, d.EntityName, d.RedisKeyPrefix, id, newEntity, updateFields, DefDemoEntityCacheExpireSeconds)
	if rerr != nil {
		d.ErrorLog([]byte("DemoSvc.UpdateById"), []byte(rerr.Error()))
	}

	return setItems != nil, merr
}

func (d *DemoSvc) ListByIds(ids ...int64) ([]*entity.DemoEntity, error) {
	var entities []*entity.DemoEntity

	err := d.SqlStore.ListByIds(d.EntityName, ids, "id desc", entity.DemoEntityType, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (d *DemoSvc) SimpleQueryAnd(sqp *define.SqlQueryParams) ([]*entity.DemoEntity, error) {
	var entities []*entity.DemoEntity

	err := d.SqlStore.SimpleQueryAnd(d.EntityName, sqp, entity.DemoEntityType, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}
