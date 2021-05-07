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

		colsValues[i] = entity.ReflectSqlInsertColValues(rev)
		ids[i] = baseEntity.Id
	}

	err := s.Dao().Insert(tableName, colNames, colsValues...).Err
	defer s.SendBackClient()

	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *SqlStore) GetById(tableName string, id int64, entityPtr interface{}) (bool, error) {
	dests := entity.ReflectSqlEntityScanDests(reflect.ValueOf(entityPtr).Elem())

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

func (s *SqlStore) UpdateById(tableName string, id int64, updateFields map[string]interface{}) (bool, error) {
	result := s.Dao().UpdateById(tableName, id, updateFields)
	defer s.SendBackClient()

	if result.Err != nil {
		return false, result.Err
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (s *SqlStore) ListByIds(tableName string, ids []int64, orderBy string, entityType reflect.Type, entityList interface{}) error {
	rows, err := s.Dao().SelectByIds(tableName, "*", orderBy, ids...)
	defer s.SendBackClient()

	if err != nil {
		return err
	}

	return entity.ReflectSqlQueryRowsToEntityList(rows, entityType, entityList)
}

func (s *SqlStore) SimpleQueryAnd(tableName string, sqp *define.SqlQueryParams, entityType reflect.Type, entityList interface{}) error {
	var orderBy string
	var offset, cnt int64
	var condItems []*mysql.SqlColQueryItem

	if sqp != nil {
		orderBy = sqp.OrderBy
		offset = sqp.Offset
		cnt = sqp.Cnt
		condItems = sqp.CondItems
	}

	rows, err := s.Dao().SimpleQueryAnd(tableName, "*", orderBy, offset, cnt, condItems...)
	defer s.SendBackClient()

	if err != nil {
		return err
	}

	return entity.ReflectSqlQueryRowsToEntityList(rows, entityType, entityList)
}

func (s *SqlStore) SimpleTotalAnd(tableName string, sqp *define.SqlQueryParams) (int64, error) {
	var condItems []*mysql.SqlColQueryItem
	if sqp != nil {
		condItems = sqp.CondItems
	}

	total, err := s.Dao().SimpleTotalAnd(tableName, condItems...)
	defer s.SendBackClient()

	return total, err
}
