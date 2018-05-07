package svc

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/mongo"

	"gdemo/dao"
	"gdemo/idgen"

	"reflect"
	"time"
)

const (
	ENTITY_MONGO_FIELD_TAG = "bson"
	ENTITY_MONGO_BASE      = "MongoBaseEntity"
)

type MongoBaseEntity struct {
	Id       int64     `bson:"_id" json:"_id"`
	AddTime  time.Time `bson:"add_time" json:"add_time"`
	EditTime time.Time `bson:"edit_time" json:"edit_time"`
}

type MongoBaseSvc struct {
	*BaseSvc

	Mclient *mongo.Client
	Mlogger golog.ILogger

	EntityName string

	Dao      *dao.MongoDao
	IdGenter *idgen.MongoIdGenter
}

type MongoQueryParams struct {
	ParamsStructPtr interface{}
	Exists          map[string]bool
	Conditions      map[string]string

	OrderBy []string
	Offset  int
	Cnt     int
}

func ReflectMongoColNames(ret reflect.Type) []string {
	var cns []string

	for i := 0; i < ret.NumField(); i++ {
		retf := ret.Field(i)
		if retf.Type.Kind() == reflect.Struct && ret.Field(i).Name == ENTITY_MONGO_BASE {
			cns = ReflectMongoColNames(retf.Type)
			continue
		}

		if name, ok := retf.Tag.Lookup(ENTITY_MONGO_FIELD_TAG); ok {
			cns = append(cns, name)
		}
	}

	return cns
}

func NewMongoBaseSvc(bs *BaseSvc, mclient *mongo.Client, entityName string) *MongoBaseSvc {

	return &MongoBaseSvc{
		BaseSvc: bs,

		Mclient:    mclient,
		EntityName: entityName,

		Dao:      dao.NewMongoDao(mclient),
		IdGenter: idgen.NewMongoIdGenter(mclient),
	}
}

func (s *MongoBaseSvc) FillBaseEntityForInsert(entity *MongoBaseEntity) error {
	id, err := s.IdGenter.GenId(s.EntityName)
	if err != nil {
		return err
	}

	ts := time.Now()
	entity.Id = id
	entity.AddTime = ts
	entity.EditTime = ts

	return nil
}

func (s *MongoBaseSvc) Insert(tableName string, colNames []string, entities ...interface{}) ([]int64, error) {
	cnt := len(entities)
	colsValues := make([][]interface{}, cnt)
	ids := make([]int64, cnt)
	for i, entity := range entities {
		rev := reflect.ValueOf(entity).Elem()
		baseEntity := rev.FieldByName(ENTITY_MONGO_BASE).Addr().Interface().(*MongoBaseEntity)
		err := s.FillBaseEntityForInsert(baseEntity)
		if err != nil {
			s.Mclient.Free()
			s.Elogger.Error([]byte("fill " + ENTITY_MONGO_BASE + " error: " + err.Error()))
			return nil, err
		}

		colsValues[i] = s.reflectInsertColValues(rev)
		ids[i] = baseEntity.Id
	}

	err := s.Dao.InsertRows(tableName, colNames, colsValues...)
	if err != nil {
		s.Elogger.Error([]byte("insert mongo error: " + err.Error()))
		return nil, err
	}

	return ids, nil
}

func (s *MongoBaseSvc) reflectInsertColValues(rev reflect.Value) []interface{} {
	var colValues []interface{}

	ret := rev.Type()
	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)

		if revf.Kind() == reflect.Struct && ret.Field(i).Name == ENTITY_MONGO_BASE {
			colValues = s.reflectInsertColValues(revf)
			continue
		}

		_, ok := ret.Field(i).Tag.Lookup(ENTITY_MONGO_FIELD_TAG)
		if ok {
			colValues = append(colValues, revf.Interface())
		}
	}

	return colValues
}

func (s *MongoBaseSvc) DeleteById(tableName string, id int64) (bool, error) {
	err := s.Dao.DeleteById(tableName, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *MongoBaseSvc) UpdateById(tableName string, id int64, newEntityPtr interface{}, updateFields map[string]bool) (error, error) {
	rnewv := reflect.ValueOf(newEntityPtr).Elem()
	oldEntity := reflect.New(rnewv.Type()).Interface()

	find, err := s.GetById(oldEntity, tableName, id)
	if err != nil {
		s.Mclient.Free()
		s.Elogger.Error([]byte("read mysql error"))
		return nil, err
	}
	if !find {
		return nil, nil
	}

	setItems := s.reflectUpdateSetItems(reflect.ValueOf(oldEntity).Elem(), rnewv, updateFields)
	if len(setItems) == 0 {
		return nil, nil
	}

	setItems["edit_time"] = time.Now()

	err = s.Dao.UpdateById(tableName, id, setItems)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *MongoBaseSvc) reflectUpdateSetItems(roldv, rnewv reflect.Value, updateFields map[string]bool) map[string]interface{} {
	setItems := make(map[string]interface{})

	rnewt := rnewv.Type()
	for i := 0; i < rnewv.NumField(); i++ {
		rnewvf := rnewv.Field(i)
		if rnewvf.Kind() == reflect.Struct && rnewt.Field(i).Name == ENTITY_MONGO_BASE {
			setItems = s.reflectUpdateSetItems(roldv.Field(i), rnewvf, updateFields)
			continue
		}

		rnewtf := rnewt.Field(i)
		colName, ok := rnewtf.Tag.Lookup(ENTITY_MONGO_FIELD_TAG)
		if !ok {
			continue
		}

		if v, ok := updateFields[colName]; !ok || !v {
			continue
		}

		setItems[colName] = rnewvf.Interface()
	}

	return setItems
}

func (s *MongoBaseSvc) GetById(entityPtr interface{}, tableName string, id int64) (bool, error) {
	result, err := s.Dao.SelectById(tableName, id)
	if err != nil {
		s.Mclient.Free()
		return false, err
	}

	err = s.Dao.ConvertJsonToStruct(result, entityPtr)
	if err != nil {
		s.Mclient.Free()
		return false, err
	}
	return true, nil
}

func (s *MongoBaseSvc) SelectAll(entityListPtr interface{}, tableName string, mqp *MongoQueryParams) error {
	setItems := s.reflectQuerySetItems(reflect.ValueOf(mqp.ParamsStructPtr).Elem(), mqp.Exists, mqp.Conditions)

	query := mongo.NewQuery().Find(setItems).Sort(mqp.OrderBy...).Skip(mqp.Offset).Limit(mqp.Cnt)
	result, err := s.Dao.SelectAll(tableName, query)

	if err != nil {
		s.Mclient.Free()
		s.Elogger.Error([]byte("list from mongo error:" + err.Error()))
		return err
	}

	err = s.Dao.ConvertJsonToStruct(result, entityListPtr)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoBaseSvc) reflectQuerySetItems(rev reflect.Value, exists map[string]bool, conditions map[string]string) map[string]interface{} {
	setItems := make(map[string]interface{})
	ret := rev.Type()

	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct {
			setItems = s.reflectQuerySetItems(revf, exists, conditions)
			continue
		}

		retf := ret.Field(i)
		name, ok := retf.Tag.Lookup(ENTITY_MONGO_FIELD_TAG)
		if !ok {
			continue
		}
		if v, ok := exists[name]; !ok || !v {
			continue
		}
		cond, ok := conditions[name]
		if !ok {
			continue
		}
		conds := make(map[string]interface{})
		conds[cond] = revf.Interface()
		setItems[name] = conds
	}

	return setItems
}
