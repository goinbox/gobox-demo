package svc

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/gomisc"
	"github.com/goinbox/mysql"

	"gdemo/idgen"

	"database/sql"
	"reflect"
	"time"
)

const (
	ENTITY_MYSQL_FIELD_TAG = "mysql"
)

var LogKindSqlBaseSvc = []byte("SqlBaseSvc")

type SqlBaseEntity struct {
	Id       int64  `mysql:"id" json:"id" redis:"id"`
	AddTime  string `mysql:"add_time" json:"add_time" redis:"add_time"`
	EditTime string `mysql:"edit_time" json:"edit_time" redis:"edit_time"`
}

type SqlBaseSvc struct {
	*BaseSvc

	Mclient *mysql.Client
	Mlogger golog.ILogger

	EntityName string

	Dao      *mysql.SqlDao
	IdGenter *idgen.SqlIdGenter
}

type SqlQueryParams struct {
	ParamsStructPtr interface{}
	Exists          map[string]bool
	Conditions      map[string]string

	OrderBy string
	Offset  int64
	Cnt     int64
}

func ReflectColNames(ret reflect.Type) []string {
	var cns []string

	for i := 0; i < ret.NumField(); i++ {
		retf := ret.Field(i)
		if retf.Type.Kind() == reflect.Struct {
			cns = ReflectColNames(retf.Type)
			continue
		}

		if name, ok := retf.Tag.Lookup(ENTITY_MYSQL_FIELD_TAG); ok {
			cns = append(cns, name)
		}
	}

	return cns
}

func NewSqlBaseSvc(bs *BaseSvc, mclient *mysql.Client, entityName string) *SqlBaseSvc {
	return &SqlBaseSvc{
		BaseSvc: bs,

		Mclient:    mclient,
		EntityName: entityName,

		Dao:      mysql.NewSqlDao(mclient),
		IdGenter: idgen.NewSqlIdGenter(mclient),
	}
}

func (s *SqlBaseSvc) FillBaseEntityForInsert(entity *SqlBaseEntity) error {
	id, err := s.IdGenter.GenId(s.EntityName)
	if err != nil {
		return err
	}

	ts := time.Now().Format(gomisc.TimeGeneralLayout())
	entity.Id = id
	entity.AddTime = ts
	entity.EditTime = ts

	return nil
}

func (s *SqlBaseSvc) Insert(tableName string, colNames []string, entities ...interface{}) ([]int64, error) {
	cnt := len(entities)
	colsValues := make([][]interface{}, cnt)
	ids := make([]int64, cnt)
	for i, entity := range entities {
		rev := reflect.ValueOf(entity).Elem()
		baseEntity := rev.FieldByName("SqlBaseEntity").Addr().Interface().(*SqlBaseEntity)
		err := s.FillBaseEntityForInsert(baseEntity)
		if err != nil {
			s.Mclient.Free()
			s.Error(LogKindSqlBaseSvc, []byte("fill SqlBaseEntity error: "+err.Error()))
			return nil, err
		}

		colsValues[i] = s.reflectInsertColValues(rev)
		ids[i] = baseEntity.Id
	}

	err := s.Dao.Insert(tableName, colNames, colsValues...).Err
	if err != nil {
		s.Error(LogKindSqlBaseSvc, []byte("insert mysql error: "+err.Error()))
		return nil, err
	}

	return ids, nil
}

func (s *SqlBaseSvc) reflectInsertColValues(rev reflect.Value) []interface{} {
	var colValues []interface{}

	ret := rev.Type()
	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct {
			colValues = s.reflectInsertColValues(revf)
			continue
		}

		_, ok := ret.Field(i).Tag.Lookup(ENTITY_MYSQL_FIELD_TAG)
		if ok {
			colValues = append(colValues, revf.Interface())
		}
	}

	return colValues
}

func (s *SqlBaseSvc) GetById(tableName string, id int64, entityPtr interface{}) (bool, error) {
	dests := s.reflectEntityScanDests(reflect.ValueOf(entityPtr).Elem())

	err := s.Dao.SelectById(tableName, "*", id).Scan(dests...)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		s.Mclient.Free()
		return false, err
	}

	return true, nil
}

func (s *SqlBaseSvc) reflectEntityScanDests(rev reflect.Value) []interface{} {
	var dests []interface{}

	ret := rev.Type()
	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct {
			dests = s.reflectEntityScanDests(revf)
			continue
		}

		_, ok := ret.Field(i).Tag.Lookup(ENTITY_MYSQL_FIELD_TAG)
		if ok {
			dests = append(dests, revf.Addr().Interface())
		}
	}

	return dests
}

func (s *SqlBaseSvc) UpdateById(tableName string, id int64, newEntityPtr interface{}, updateFields map[string]bool) ([]*mysql.SqlColQueryItem, error) {
	rnewv := reflect.ValueOf(newEntityPtr).Elem()
	oldEntity := reflect.New(rnewv.Type()).Interface()

	find, err := s.GetById(tableName, id, oldEntity)
	if err != nil {
		s.Mclient.Free()
		s.Error(LogKindSqlBaseSvc, []byte("read mysql error"))
		return nil, err
	}
	if !find {
		return nil, nil
	}

	setItems := s.reflectUpdateSetItems(reflect.ValueOf(oldEntity).Elem(), rnewv, updateFields)
	if len(setItems) == 0 {
		return nil, nil
	}

	setItems = append(setItems, mysql.NewSqlColQueryItem("edit_time", "", time.Now().Format(gomisc.TimeGeneralLayout())))
	result := s.Dao.UpdateById(tableName, id, setItems...)
	if result.Err != nil {
		s.Error(LogKindSqlBaseSvc, []byte("update mysql error: "+result.Err.Error()))
		return nil, result.Err
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return setItems, nil
}

func (s *SqlBaseSvc) reflectUpdateSetItems(roldv, rnewv reflect.Value, updateFields map[string]bool) []*mysql.SqlColQueryItem {
	var setItems []*mysql.SqlColQueryItem

	rnewt := rnewv.Type()
	for i := 0; i < rnewv.NumField(); i++ {
		rnewvf := rnewv.Field(i)
		if rnewvf.Kind() == reflect.Struct {
			setItems = s.reflectUpdateSetItems(roldv.Field(i), rnewvf, updateFields)
			continue
		}

		rnewtf := rnewt.Field(i)
		colName, ok := rnewtf.Tag.Lookup(ENTITY_MYSQL_FIELD_TAG)
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

func (s *SqlBaseSvc) ListByIds(tableName string, ids []int64, orderBy string, entityType reflect.Type, listPtr interface{}) error {
	rows, err := s.Dao.SelectByIds(tableName, "*", orderBy, ids...)
	if err != nil {
		s.Mclient.Free()
		s.Error(LogKindSqlBaseSvc, []byte("list from mysql error:"+err.Error()))
		return err
	}

	err = s.ReflectQueryRowsToEntityList(rows, entityType, listPtr)
	if err != nil {
		s.Error(LogKindSqlBaseSvc, []byte("list from mysql error:"+err.Error()))
		return err
	}

	return nil
}

func (s *SqlBaseSvc) ReflectQueryRowsToEntityList(rows *sql.Rows, ret reflect.Type, listPtr interface{}) error {
	if rows.Next() == false {
		return nil
	}

	rlistv := reflect.ValueOf(listPtr).Elem()
	rev := reflect.New(ret)
	dests := s.reflectEntityScanDests(rev.Elem())
	err := rows.Scan(dests...)
	if err != nil {
		s.Error(LogKindSqlBaseSvc, []byte("list from mysql error:"+err.Error()))
		return err
	}
	rlistv.Set(reflect.Append(rlistv, rev))

	for rows.Next() {
		rev = reflect.New(ret)
		dests = s.reflectEntityScanDests(rev.Elem())
		err = rows.Scan(dests...)
		if err != nil {
			s.Error(LogKindSqlBaseSvc, []byte("list from mysql error:"+err.Error()))
			return err
		}
		rlistv.Set(reflect.Append(rlistv, rev))
	}

	return nil
}

func (s *SqlBaseSvc) SimpleQueryAnd(tableName string, sqp *SqlQueryParams, entityType reflect.Type, listPtr interface{}) error {
	setItems := s.reflectQuerySetItems(reflect.ValueOf(sqp.ParamsStructPtr).Elem(), sqp.Exists, sqp.Conditions)

	rows, err := s.Dao.SimpleQueryAnd(tableName, "*", sqp.OrderBy, sqp.Offset, sqp.Cnt, setItems...)
	if err != nil {
		s.Mclient.Free()
		s.Error(LogKindSqlBaseSvc, []byte("list from mysql error:"+err.Error()))
		return err
	}

	err = s.ReflectQueryRowsToEntityList(rows, entityType, listPtr)
	if err != nil {
		s.Error(LogKindSqlBaseSvc, []byte("list from mysql error:"+err.Error()))
		return err
	}

	return nil
}

func (s *SqlBaseSvc) SimpleTotalAnd(tableName string, sqp *SqlQueryParams) (int64, error) {
	setItems := s.reflectQuerySetItems(reflect.ValueOf(sqp.ParamsStructPtr).Elem(), sqp.Exists, sqp.Conditions)

	return s.Dao.SimpleTotalAnd(tableName, setItems...)
}

func (s *SqlBaseSvc) reflectQuerySetItems(rev reflect.Value, exists map[string]bool, conditions map[string]string) []*mysql.SqlColQueryItem {
	var setItems []*mysql.SqlColQueryItem
	ret := rev.Type()

	for i := 0; i < rev.NumField(); i++ {
		revf := rev.Field(i)
		if revf.Kind() == reflect.Struct {
			setItems = s.reflectQuerySetItems(revf, exists, conditions)
			continue
		}

		retf := ret.Field(i)
		name, ok := retf.Tag.Lookup(ENTITY_MYSQL_FIELD_TAG)
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
