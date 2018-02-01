package svc

import (
	"github.com/goinbox/mysql"
	"github.com/goinbox/redis"
	"github.com/goinbox/golog"

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

func (this *DemoSvc) Insert(entities ...*DemoEntity) ([]int64, error) {
	is := make([]interface{}, len(entities))
	for i, entity := range entities {
		is[i] = entity
	}

	return this.SqlRedisBindSvc.Insert(this.entityName, demoColNames, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS, is...)
}

func (this *DemoSvc) GetById(id int64) (*DemoEntity, error) {
	entity := new(DemoEntity)

	find, err := this.SqlRedisBindSvc.GetById(this.entityName, id, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS, entity)
	if err != nil {
		return nil, err
	}
	if !find {
		return nil, nil
	}

	return entity, nil
}

func (this *DemoSvc) DeleteById(id int64) (bool, error) {
	return this.SqlRedisBindSvc.DeleteById(this.entityName, id)
}

func (this *DemoSvc) UpdateById(id int64, newEntity *DemoEntity, updateFields map[string]bool) (bool, error) {
	return this.SqlRedisBindSvc.UpdateById(this.entityName, id, newEntity, updateFields, DEF_DEMO_ENTITY_CACHE_EXPIRE_SECONDS)
}

func (this *DemoSvc) ListByIds(ids ...int64) ([]*DemoEntity, error) {
	var entities []*DemoEntity

	err := this.SqlBaseSvc.ListByIds(this.entityName, ids, "id desc", demoEntityType, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (this *DemoSvc) SimpleQueryAnd(sqp *SqlQueryParams) ([]*DemoEntity, error) {
	var entities []*DemoEntity

	err := this.SqlBaseSvc.SimpleQueryAnd(this.entityName, sqp, demoEntityType, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}
