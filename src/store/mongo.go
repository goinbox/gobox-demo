package store

import (
	"reflect"
	"time"

	"github.com/globalsign/mgo"
	"github.com/goinbox/mongo"

	"gdemo/define"
	"gdemo/define/entity"
	"gdemo/idgen"
	"gdemo/resource"
)

type MongoStore struct {
	traceId  []byte
	pool     *mongo.Pool
	useIdGen bool

	dao      *mongo.MongoDao
	idGenter *idgen.MongoIdGenter
}

func NewMongoStore(traceId []byte, pool *mongo.Pool, useIdGen bool) *MongoStore {
	return &MongoStore{
		traceId:  traceId,
		pool:     pool,
		useIdGen: useIdGen,
	}
}

func (s *MongoStore) Dao() *mongo.MongoDao {
	if s.dao == nil {
		s.dao = &mongo.MongoDao{}
	}

	if s.dao.Client == nil {
		s.dao.Client, _ = s.pool.Get()
		s.dao.Client.SetLogger(resource.AccessLogger).SetTraceId(s.traceId)
	}

	return s.dao
}

func (s *MongoStore) IdGenter() *idgen.MongoIdGenter {
	if !s.useIdGen {
		return nil
	}

	if s.idGenter == nil {
		s.idGenter = idgen.NewMongoIdGenter(s.Dao().Client)
	}

	return s.idGenter
}

func (s *MongoStore) SendBackClient() {
	if !s.dao.Client.Closed() {
		s.dao.Client.SetLogger(resource.NoopLogger)
		_ = s.pool.Put(s.dao.Client)
	}

	s.dao.Client = nil
	if s.idGenter != nil {
		s.idGenter.SetClient(nil)
	}
}

func (s *MongoStore) FillBaseEntityForInsert(baseEntity *entity.MongoBaseEntity, rev reflect.Value, tableName string) error {
	ts := time.Now()
	entityId := rev.FieldByName("Id").Interface()
	entityAddTime := rev.FieldByName("AddTime").Interface().(time.Time)
	entityEditTime := rev.FieldByName("EditTime").Interface().(time.Time)
	if entityId == nil {
		id, err := s.IdGenter().GenId(tableName)
		if err != nil {
			return err
		}
		baseEntity.Id = id
	} else {
		baseEntity.Id = entityId
	}
	if entityAddTime.IsZero() {
		baseEntity.AddTime = ts
	} else {
		baseEntity.AddTime = entityAddTime
	}
	if entityEditTime.IsZero() {
		baseEntity.EditTime = ts
	} else {
		baseEntity.EditTime = entityEditTime
	}
	return nil
}

func (s *MongoStore) Insert(tableName string, colNames []string, entities ...interface{}) ([]interface{}, error) {
	cnt := len(entities)
	colsValues := make([][]interface{}, cnt)
	ids := make([]interface{}, cnt)
	for i, item := range entities {
		rev := reflect.ValueOf(item).Elem()
		if rev.FieldByName(entity.EntityMongoBase).IsValid() {
			baseEntity := rev.FieldByName(entity.EntityMongoBase).Addr().Interface().(*entity.MongoBaseEntity)
			err := s.FillBaseEntityForInsert(baseEntity, rev, tableName)
			if err != nil {
				return nil, err
			}

			colsValues[i] = s.reflectInsertColValues(rev)
			ids[i] = baseEntity.Id
		} else {
			colsValues[i] = s.reflectInsertColValues(rev)
			ids[i] = rev.FieldByName("Id").Interface()
		}
	}

	err := s.Dao().InsertRows(tableName, colNames, colsValues...)
	defer s.SendBackClient()
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *MongoStore) reflectInsertColValues(rev reflect.Value) []interface{} {
	var colValues []interface{}

	ret := rev.Type()
	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct && revf.Type().Name() == entity.EntityMongoBase {
			colValues = s.reflectInsertColValues(revf)
			continue
		}

		_, ok := ret.Field(i).Tag.Lookup(entity.EntityMongoFieldTag)
		if ok {
			colValues = append(colValues, revf.Interface())
		}
	}

	return colValues
}

func (s *MongoStore) DeleteById(tableName string, id interface{}) (bool, error) {
	err := s.Dao().DeleteById(tableName, id)
	defer s.SendBackClient()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *MongoStore) UpdateById(tableName string, id interface{}, newEntityPtr interface{}, updateFields map[string]bool) (map[string]interface{}, error) {
	rnewv := reflect.ValueOf(newEntityPtr).Elem()
	oldEntity := reflect.New(rnewv.Type()).Interface()

	find, err := s.GetById(oldEntity, tableName, id)
	if err != nil {
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

	err = s.Dao().UpdateById(tableName, id, setItems)
	defer s.SendBackClient()
	if err != nil {
		return nil, err
	}

	return setItems, nil
}

func (s *MongoStore) reflectUpdateSetItems(roldv, rnewv reflect.Value, updateFields map[string]bool) map[string]interface{} {
	setItems := make(map[string]interface{})

	rnewt := rnewv.Type()
	for i := 0; i < rnewv.NumField(); i++ {
		rnewvf := rnewv.Field(i)
		if rnewvf.Kind() == reflect.Struct && rnewvf.Type().Name() == entity.EntityMongoBase {
			setItems = s.reflectUpdateSetItems(roldv.Field(i), rnewvf, updateFields)
			continue
		}

		rnewtf := rnewt.Field(i)
		colName, ok := rnewtf.Tag.Lookup(entity.EntityMongoFieldTag)
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

func (s *MongoStore) GetById(entityPtr interface{}, tableName string, id interface{}) (bool, error) {
	result, err := s.Dao().SelectById(tableName, id)
	defer s.SendBackClient()
	if err != nil && err != mgo.ErrNotFound {
		return false, err
	}

	if err == mgo.ErrNotFound {
		return false, nil
	}

	err = s.Dao().ConvertJsonToStruct(result, entityPtr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *MongoStore) SelectAll(entityListPtr interface{}, tableName string, mqp *define.MongoQueryParams, setItems map[string]interface{}) error {
	if setItems == nil {
		setItems = s.ReflectQuerySetItems(reflect.ValueOf(mqp.ParamsStructPtr).Elem(), mqp.Exists, mqp.Conditions)
	}

	query := mongo.NewQuery().Find(setItems).Sort(mqp.OrderBy...).Skip(mqp.Offset).Limit(mqp.Cnt)
	result, err := s.Dao().SelectAll(tableName, query)
	defer s.SendBackClient()

	if err != nil {
		return err
	}

	err = s.Dao().ConvertJsonToStruct(result, entityListPtr)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoStore) ReflectQuerySetItems(rev reflect.Value, exists map[string]bool, conditions map[string]string) map[string]interface{} {
	setItems := make(map[string]interface{})
	ret := rev.Type()

	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct && revf.Type().Name() == entity.EntityMongoBase {
			setItems = s.ReflectQuerySetItems(revf, exists, conditions)
			continue
		}

		retf := ret.Field(i)
		name, ok := retf.Tag.Lookup(entity.EntityMongoFieldTag)
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
