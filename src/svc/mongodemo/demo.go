package mongodemo

import (
	"reflect"

	"github.com/goinbox/mongo"

	"gdemo/define"
	"gdemo/define/entity"
	"gdemo/resource"
	"gdemo/store"
	"gdemo/svc"
)

type MongoDemoSvc struct {
	*svc.BaseSvc
	*store.MongoStore
	EntityName string
}

func NewMongoDemoSvc(traceId []byte) *MongoDemoSvc {
	ms := store.NewMongoStore(traceId, resource.MongoClientPool, true)

	return &MongoDemoSvc{
		BaseSvc: &svc.BaseSvc{
			TraceId: traceId,
		},
		MongoStore: ms,
		EntityName: "demo",
	}
}

func (d *MongoDemoSvc) Insert(entities ...*entity.MongoDemoEntity) ([]interface{}, error) {
	is := make([]interface{}, len(entities))
	for i, item := range entities {
		is[i] = item
	}
	ids, err := d.MongoStore.Insert(d.EntityName, entity.MongoDemoColNames, is...)
	if err != nil {
		d.ErrorLog([]byte("MongoDemoSvc.Insert"), []byte(err.Error()))
	}
	return ids, err
}

func (d *MongoDemoSvc) DeleteById(id interface{}) (bool, error) {
	find, err := d.MongoStore.DeleteById(d.EntityName, id)
	if err != nil {
		d.ErrorLog([]byte("MongoDemoSvc.DeleteById"), []byte(err.Error()))
	}
	return find, err
}

func (d *MongoDemoSvc) UpdateById(id interface{}, newEntity *entity.MongoDemoEntity, updateFields map[string]bool) (bool, error) {
	setItems, err := d.MongoStore.UpdateById(d.EntityName, id, newEntity, updateFields)
	if err != nil {
		d.ErrorLog([]byte("MongoDemoSvc.UpdateById"), []byte(err.Error()))
	}
	return (setItems != nil), err
}

func (d *MongoDemoSvc) GetById(id interface{}) (*entity.MongoDemoEntity, error) {
	entity := new(entity.MongoDemoEntity)
	find, err := d.MongoStore.GetById(entity, d.EntityName, id)
	if err != nil {
		d.ErrorLog([]byte("MongoDemoSvc.GetById"), []byte(err.Error()))
	}
	if err != nil {
		return nil, err
	}
	if !find {
		return nil, nil
	}

	return entity, nil
}

func (d *MongoDemoSvc) SelectAll(mqp *define.MongoQueryParams) (*[]entity.MongoDemoEntity, error) {
	entities := new([]entity.MongoDemoEntity)
	err := d.MongoStore.SelectAll(entities, d.EntityName, mqp, nil)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (d *MongoDemoSvc) SelectRegex(mqp *define.MongoQueryParams) (*[]entity.MongoDemoEntity, error) {
	entities := new([]entity.MongoDemoEntity)

	setItems := d.MongoStore.ReflectQuerySetItems(reflect.ValueOf(mqp.ParamsStructPtr).Elem(), mqp.Exists, mqp.Conditions)
	setItems["name"].(map[string]interface{})[mongo.MONGO_COND_OPTIONS] = "i"

	err := d.MongoStore.SelectAll(entities, d.EntityName, mqp, setItems)
	if err != nil {
		return nil, err
	}
	return entities, nil
}
