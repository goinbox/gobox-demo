package mongodemo

import (
	"gdemo/resource"
	"gdemo/svc"

	"github.com/goinbox/mongo"

	"reflect"
)

var MongoDemoEntityType reflect.Type = reflect.TypeOf(MongoDemoEntity{})
var demoColNames []string = svc.ReflectMongoColNames(MongoDemoEntityType)

type MongoDemoEntity struct {
	svc.MongoBaseEntity

	Name   string `bson:"name" json:"name"`
	Status int    `bson:"status" json:"status"`
}

type MongoDemoSvc struct {
	*svc.BaseSvc
	*svc.MongoSvc
	EntityName string
}

func NewMongoDemoSvc(traceId []byte) *MongoDemoSvc {
	ms := svc.NewMongoSvc(traceId, resource.MongoClientPool, true)

	return &MongoDemoSvc{
		BaseSvc: &svc.BaseSvc{
			TraceId: traceId,
		},
		MongoSvc:   ms,
		EntityName: "demo",
	}
}

func (d *MongoDemoSvc) Insert(entities ...*MongoDemoEntity) ([]interface{}, error) {
	is := make([]interface{}, len(entities))
	for i, entity := range entities {
		is[i] = entity
	}

	return d.MongoSvc.Insert(d.EntityName, demoColNames, is...)
}

func (d *MongoDemoSvc) DeleteById(id interface{}) (bool, error) {
	return d.MongoSvc.DeleteById(d.EntityName, id)
}

func (d *MongoDemoSvc) UpdateById(id interface{}, newEntity *MongoDemoEntity, updateFields map[string]bool) error {
	return d.MongoSvc.UpdateById(d.EntityName, id, newEntity, updateFields)
}

func (d *MongoDemoSvc) GetById(id interface{}) (*MongoDemoEntity, error) {
	entity := new(MongoDemoEntity)
	find, err := d.MongoSvc.GetById(entity, d.EntityName, id)
	if err != nil {
		return nil, err
	}
	if !find {
		return nil, nil
	}

	return entity, nil
}

func (d *MongoDemoSvc) SelectAll(mqp *svc.MongoQueryParams) (*[]MongoDemoEntity, error) {
	entities := new([]MongoDemoEntity)
	err := d.MongoSvc.SelectAll(entities, d.EntityName, mqp, nil)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (d *MongoDemoSvc) SelectRegex(mqp *svc.MongoQueryParams) (*[]MongoDemoEntity, error) {
	entities := new([]MongoDemoEntity)

	setItems := d.MongoSvc.ReflectQuerySetItems(reflect.ValueOf(mqp.ParamsStructPtr).Elem(), mqp.Exists, mqp.Conditions)
	setItems["name"].(map[string]interface{})[mongo.MONGO_COND_OPTIONS] = "i"

	err := d.MongoSvc.SelectAll(entities, d.EntityName, mqp, setItems)
	if err != nil {
		return nil, err
	}
	return entities, nil
}
