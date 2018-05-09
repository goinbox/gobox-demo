package mongodemo

import (
	"gdemo/dao"
	"gdemo/svc"

	"github.com/goinbox/golog"
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
	*svc.MongoBaseSvc
}

func NewMongoDemoSvc(elogger golog.ILogger, mclient *mongo.Client) *MongoDemoSvc {
	bs := svc.NewBaseSvc(elogger)
	sbs := svc.NewMongoBaseSvc(bs, mclient, "mycoll")

	return &MongoDemoSvc{
		sbs,
	}
}

func (d *MongoDemoSvc) Insert(entities ...*MongoDemoEntity) ([]interface{}, error) {
	is := make([]interface{}, len(entities))
	for i, entity := range entities {
		is[i] = entity
	}

	return d.MongoBaseSvc.Insert(d.EntityName, demoColNames, is...)
}

func (d *MongoDemoSvc) DeleteById(id interface{}) (bool, error) {
	return d.MongoBaseSvc.DeleteById(d.EntityName, id)
}

func (d *MongoDemoSvc) UpdateById(id interface{}, newEntity *MongoDemoEntity, updateFields map[string]bool) (error, error) {
	return d.MongoBaseSvc.UpdateById(d.EntityName, id, newEntity, updateFields)
}

func (d *MongoDemoSvc) GetById(id interface{}) (*MongoDemoEntity, error) {
	entity := new(MongoDemoEntity)
	find, err := d.MongoBaseSvc.GetById(entity, d.EntityName, id)
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
	err := d.MongoBaseSvc.SelectAll(entities, d.EntityName, mqp, nil)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (d *MongoDemoSvc) SelectRegex(mqp *svc.MongoQueryParams) (*[]MongoDemoEntity, error) {
	entities := new([]MongoDemoEntity)

	setItems := d.MongoBaseSvc.ReflectQuerySetItems(reflect.ValueOf(mqp.ParamsStructPtr).Elem(), mqp.Exists, mqp.Conditions)
	setItems["name"].(map[string]interface{})[dao.MONGO_COND_OPTIONS] = "i"

	err := d.MongoBaseSvc.SelectAll(entities, d.EntityName, mqp, setItems)
	if err != nil {
		return nil, err
	}
	return entities, nil
}
