package demo

import (
	"gdemo/svc"

	"github.com/goinbox/golog"
	"github.com/goinbox/mysql"
	"github.com/goinbox/redis"

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
	*svc.SqlRedisBindSvc
}

func NewDemoSvc(elogger golog.ILogger, mclient *mysql.Client, redisKeyPrefix string, rclient *redis.Client) *DemoSvc {
	bs := svc.NewBaseSvc(elogger)
	sbs := svc.NewSqlBaseSvc(bs, mclient, "demo")

	return &DemoSvc{
		&svc.SqlRedisBindSvc{
			BaseSvc:      bs,
			SqlBaseSvc:   sbs,
			RedisBaseSvc: svc.NewRedisBaseSvc(bs, rclient),

			RedisKeyPrefix: redisKeyPrefix,
		},
	}
}

func (d *DemoSvc) Insert(entities ...*DemoEntity) ([]int64, error) {
	is := make([]interface{}, len(entities))
	for i, entity := range entities {
		is[i] = entity
	}

	return d.SqlRedisBindSvc.Insert(d.EntityName, demoColNames, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS, is...)
}

func (d *DemoSvc) GetById(id int64) (*DemoEntity, error) {
	entity := new(DemoEntity)

	find, err := d.SqlRedisBindSvc.GetById(d.EntityName, id, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS, entity)
	if err != nil {
		return nil, err
	}
	if !find {
		return nil, nil
	}

	return entity, nil
}

func (d *DemoSvc) DeleteById(id int64) (bool, error) {
	return d.SqlRedisBindSvc.DeleteById(d.EntityName, id)
}

func (d *DemoSvc) UpdateById(id int64, newEntity *DemoEntity, updateFields map[string]bool) (bool, error) {
	return d.SqlRedisBindSvc.UpdateById(d.EntityName, id, newEntity, updateFields, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS)
}

func (d *DemoSvc) ListByIds(ids ...int64) ([]*DemoEntity, error) {
	var entities []*DemoEntity

	err := d.SqlBaseSvc.ListByIds(d.EntityName, ids, "id desc", demoEntityType, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (d *DemoSvc) SimpleQueryAnd(sqp *svc.SqlQueryParams) ([]*DemoEntity, error) {
	var entities []*DemoEntity

	err := d.SqlBaseSvc.SimpleQueryAnd(d.EntityName, sqp, demoEntityType, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}
