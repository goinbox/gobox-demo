package store

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/goinbox/gomisc"
	"github.com/goinbox/mysql"

	"gdemo/define"
	"gdemo/define/entity"
	"gdemo/idgen"
	"gdemo/resource"
)

type SqlStore struct {
	traceId  []byte
	pool     *mysql.Pool
	useIdGen bool

	dao      *mysql.SqlDao
	idGenter *idgen.SqlIdGenter
}

func NewSqlStore(traceId []byte, pool *mysql.Pool, useIdGen bool) *SqlStore {
	return &SqlStore{
		traceId:  traceId,
		pool:     pool,
		useIdGen: useIdGen,
	}
}

func (s *SqlStore) Dao() *mysql.SqlDao {
	if s.dao == nil {
		s.dao = &mysql.SqlDao{}
	}

	if s.dao.Client == nil {
		s.dao.Client, _ = s.pool.Get()
		s.dao.Client.SetLogger(resource.AccessLogger).SetTraceId(s.traceId)
	}

	return s.dao
}

func (s *SqlStore) IdGenter() *idgen.SqlIdGenter {
	if !s.useIdGen {
		return nil
	}

	if s.idGenter == nil {
		s.idGenter = idgen.NewSqlIdGenter(s.Dao().Client)
	}

	return s.idGenter
}

func (s *SqlStore) SendBackClient() {
	if !s.dao.Client.Closed() {
		s.dao.Client.SetLogger(resource.NoopLogger)
		_ = s.pool.Put(s.dao.Client)
	}

	s.dao.Client = nil
	if s.idGenter != nil {
		s.idGenter.SetClient(nil)
	}
}

func (s *SqlStore) Renew(traceId []byte, pool *mysql.Pool, useIdGen bool) *SqlStore {
	if s.dao != nil && s.dao.Client != nil {
		s.SendBackClient()
	}

	s.traceId = traceId
	s.pool = pool

	return s
}

func (s *SqlStore) SetPool(pool *mysql.Pool) *SqlStore {
	return s.Renew(s.traceId, pool, s.useIdGen)
}

func (s *SqlStore) FillBaseEntityForInsert(entity *entity.SqlBaseEntity, name string) error {
	if s.useIdGen {
		id, err := s.IdGenter().GenId(name)
		if err != nil {
			return err
		}
		entity.Id = id
	}

	ts := time.Now().Format(gomisc.TimeGeneralLayout())
	entity.AddTime = ts
	entity.EditTime = ts

	return nil
}

func (s *SqlStore) Insert(tableName, entityName string, colNames []string, entities ...interface{}) ([]int64, error) {
	cnt := len(entities)
	colsValues := make([][]interface{}, cnt)
	ids := make([]int64, cnt)
	for i, ent := range entities {
		rev := reflect.ValueOf(ent).Elem()
		baseEntity := rev.FieldByName("SqlBaseEntity").Addr().Interface().(*entity.SqlBaseEntity)
		err := s.FillBaseEntityForInsert(baseEntity, entityName)
		if err != nil {
			return nil, err
		}

		colsValues[i] = s.ReflectInsertColValues(rev)
		ids[i] = baseEntity.Id
	}

	err := s.Dao().Insert(tableName, colNames, colsValues...).Err
	defer s.SendBackClient()

	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *SqlStore) ReflectInsertColValues(rev reflect.Value) []interface{} {
	var colValues []interface{}

	ret := rev.Type()
	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct {
			colValues = s.ReflectInsertColValues(revf)
			continue
		}

		_, ok := ret.Field(i).Tag.Lookup(entity.EntityMysqlFieldTag)
		if ok {
			colValues = append(colValues, revf.Interface())
		}
	}

	return colValues
}

func (s *SqlStore) GetById(tableName string, id int64, entityPtr interface{}) (bool, error) {
	dests := s.ReflectEntityScanDests(reflect.ValueOf(entityPtr).Elem())

	err := s.Dao().SelectById(tableName, "*", id).Scan(dests...)
	defer s.SendBackClient()

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *SqlStore) ReflectEntityScanDests(rev reflect.Value) []interface{} {
	var dests []interface{}

	ret := rev.Type()
	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct {
			dests = s.ReflectEntityScanDests(revf)
			continue
		}

		_, ok := ret.Field(i).Tag.Lookup(entity.EntityMysqlFieldTag)
		if ok {
			dests = append(dests, revf.Addr().Interface())
		}
	}

	return dests
}

func (s *SqlStore) UpdateById(tableName string, id int64, newEntityPtr interface{}, updateFields map[string]bool) ([]*mysql.SqlColQueryItem, error) {
	rnewv := reflect.ValueOf(newEntityPtr).Elem()
	oldEntity := reflect.New(rnewv.Type()).Interface()

	find, err := s.GetById(tableName, id, oldEntity)
	if err != nil {
		return nil, err
	}
	if !find {
		return nil, nil
	}

	setItems := s.ReflectUpdateSetItems(reflect.ValueOf(oldEntity).Elem(), rnewv, updateFields)
	if len(setItems) == 0 {
		return nil, nil
	}

	setItems = append(setItems, mysql.NewSqlColQueryItem("edit_time", "", time.Now().Format(gomisc.TimeGeneralLayout())))
	result := s.Dao().UpdateById(tableName, id, setItems...)
	defer s.SendBackClient()

	if result.Err != nil {
		return nil, result.Err
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return setItems, nil
}

func (s *SqlStore) ReflectUpdateSetItems(roldv, rnewv reflect.Value, updateFields map[string]bool) []*mysql.SqlColQueryItem {
	var setItems []*mysql.SqlColQueryItem

	rnewt := rnewv.Type()
	for i := 0; i < rnewv.NumField(); i++ {
		rnewvf := rnewv.Field(i)
		if rnewvf.Kind() == reflect.Struct {
			setItems = s.ReflectUpdateSetItems(roldv.Field(i), rnewvf, updateFields)
			continue
		}

		rnewtf := rnewt.Field(i)
		colName, ok := rnewtf.Tag.Lookup(entity.EntityMysqlFieldTag)
		if !ok {
			continue
		}
		if v, ok := updateFields[colName]; !ok || !v {
			continue
		}

		nv := rnewvf.Interface()
		if nv != roldv.Field(i).Interface() {
			setItems = append(setItems, mysql.NewSqlColQueryItem(colName, "", nv))
		}
	}

	return setItems
}

func (s *SqlStore) ListByIds(tableName string, ids []int64, orderBy string, entityType reflect.Type, listPtr interface{}) error {
	rows, err := s.Dao().SelectByIds(tableName, "*", orderBy, ids...)
	defer s.SendBackClient()

	if err != nil {
		return err
	}

	return s.ReflectQueryRowsToEntityList(rows, entityType, listPtr)
}

func (s *SqlStore) ReflectQueryRowsToEntityList(rows *sql.Rows, ret reflect.Type, listPtr interface{}) error {
	if rows.Next() == false {
		return nil
	}

	rlistv := reflect.ValueOf(listPtr).Elem()
	rev := reflect.New(ret)
	dests := s.ReflectEntityScanDests(rev.Elem())
	err := rows.Scan(dests...)
	if err != nil {
		return err
	}
	rlistv.Set(reflect.Append(rlistv, rev))

	for rows.Next() {
		rev = reflect.New(ret)
		dests = s.ReflectEntityScanDests(rev.Elem())
		err = rows.Scan(dests...)
		if err != nil {
			return err
		}
		rlistv.Set(reflect.Append(rlistv, rev))
	}

	return nil
}

func (s *SqlStore) SimpleQueryAnd(tableName string, sqp *define.SqlQueryParams, entityType reflect.Type, listPtr interface{}) error {
	var orderBy string
	var offset, cnt int64
	var setItems []*mysql.SqlColQueryItem

	if sqp != nil {
		orderBy = sqp.OrderBy
		offset = sqp.Offset
		cnt = sqp.Cnt
		if sqp.ParamsStructPtr != nil {
			setItems = s.ReflectQuerySetItems(reflect.ValueOf(sqp.ParamsStructPtr).Elem(), sqp.Exists, sqp.Conditions)
		}
	}

	rows, err := s.Dao().SimpleQueryAnd(tableName, "*", orderBy, offset, cnt, setItems...)
	defer s.SendBackClient()

	if err != nil {
		return err
	}

	return s.ReflectQueryRowsToEntityList(rows, entityType, listPtr)
}

func (s *SqlStore) SimpleTotalAnd(tableName string, sqp *define.SqlQueryParams) (int64, error) {
	var setItems []*mysql.SqlColQueryItem
	if sqp != nil && sqp.ParamsStructPtr != nil {
		setItems = s.ReflectQuerySetItems(reflect.ValueOf(sqp.ParamsStructPtr).Elem(), sqp.Exists, sqp.Conditions)
	}

	total, err := s.Dao().SimpleTotalAnd(tableName, setItems...)
	defer s.SendBackClient()

	return total, err
}

func (s *SqlStore) ReflectQuerySetItems(rev reflect.Value, exists map[string]bool, conditions map[string]string) []*mysql.SqlColQueryItem {
	var setItems []*mysql.SqlColQueryItem
	ret := rev.Type()

	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct {
			setItems = s.ReflectQuerySetItems(revf, exists, conditions)
			continue
		}

		retf := ret.Field(i)
		name, ok := retf.Tag.Lookup(entity.EntityMysqlFieldTag)
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

		setItems = append(setItems, mysql.NewSqlColQueryItem(name, cond, revf.Interface()))
	}

	return setItems
}
