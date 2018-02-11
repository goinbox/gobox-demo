package svc

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/mysql"
	"github.com/goinbox/redis"

	"reflect"
)

const (
	DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS = 60 * 30
)

var demoEntityType reflect.Type = reflect.TypeOf(DemoEntity{})
var demoColNames []string = ReflectColNames(demoEntityType)

type DemoEntity struct {
	SqlBaseEntity

	Name   string `mysql:"name" json:"name" redis:"name"`
	Status int    `mysql:"status" json:"status" redis:"status"`
}

type DemoSvc struct {
	*SqlRedisBindSvc
}

func NewDemoSvc(elogger golog.ILogger, mclient *mysql.Client, redisKeyPrefix string, rclient *redis.Client) *DemoSvc {
	bs := NewBaseSvc(elogger)
	sbs := NewSqlBaseSvc(bs, mclient, "demo")

	return &DemoSvc{
		&SqlRedisBindSvc{
			BaseSvc:      bs,
			SqlBaseSvc:   sbs,
			RedisBaseSvc: NewRedisBaseSvc(bs, rclient),

			redisKeyPrefix: redisKeyPrefix,
		},
	}
}

func (d *DemoSvc) Insert(entities ...*DemoEntity) ([]int64, error) {
	is := make([]interface{}, len(entities))
	for i, entity := range entities {
		is[i] = entity
	}

	return d.SqlRedisBindSvc.Insert(d.entityName, demoColNames, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS, is...)
}

func (d *DemoSvc) GetById(id int64) (*DemoEntity, error) {
	entity := new(DemoEntity)

	find, err := d.SqlRedisBindSvc.GetById(d.entityName, id, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS, entity)
	if err != nil {
		return nil, err
	}
	if !find {
		return nil, nil
	}

	return entity, nil
}

func (d *DemoSvc) DeleteById(id int64) (bool, error) {
	return d.SqlRedisBindSvc.DeleteById(d.entityName, id)
}

func (d *DemoSvc) UpdateById(id int64, newEntity *DemoEntity, updateFields map[string]bool) (bool, error) {
	return d.SqlRedisBindSvc.UpdateById(d.entityName, id, newEntity, updateFields, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS)
}

func (d *DemoSvc) ListByIds(ids ...int64) ([]*DemoEntity, error) {
	var entities []*DemoEntity

	err := d.SqlBaseSvc.ListByIds(d.entityName, ids, "id desc", demoEntityType, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (d *DemoSvc) SimpleQueryAnd(sqp *SqlQueryParams) ([]*DemoEntity, error) {
	var entities []*DemoEntity

	err := d.SqlBaseSvc.SimpleQueryAnd(d.entityName, sqp, demoEntityType, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}
