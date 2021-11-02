package model

import "github.com/goinbox/mysql"

const (
	ColumnNameID = "id"
)

type Dao interface {
	Insert(entities ...interface{}) error

	DeleteByIDs(ids ...int64) error
	UpdateByIDs(fields map[string]interface{}, ids ...int64) error
	SelectByIDs(dest interface{}, ids ...int64) error

	SimpleQueryAnd(params *mysql.SqlQueryParams, dest interface{}) error
	SimpleTotalAnd(items ...*mysql.SqlColQueryItem) (int64, error)

	Begin() error
	Commit() error
	Rollback() error
}

type BaseDao struct {
	tableName string
	dao       *mysql.EntityDao
}

func NewBaseDao(tableName string, client *mysql.Client) *BaseDao {
	return &BaseDao{
		tableName: tableName,
		dao: &mysql.EntityDao{
			Dao: mysql.Dao{
				Client: client,
			},
		},
	}
}

func (d *BaseDao) Insert(entities ...interface{}) error {
	return d.dao.InsertEntities(d.tableName, entities...)
}

func (d *BaseDao) DeleteByIDs(ids ...int64) error {
	if len(ids) == 1 {
		return d.dao.DeleteById(d.tableName, ids[0]).Err
	}

	sqb := new(mysql.SqlQueryBuilder)
	sqb.Delete(d.tableName).
		WhereConditionAnd(&mysql.SqlColQueryItem{
			Name:      ColumnNameID,
			Condition: mysql.SqlCondIn,
			Value:     ids,
		})

	_, err := d.dao.Exec(sqb.Query(), sqb.Args()...)
	return err
}

func (d *BaseDao) UpdateByIDs(fields map[string]interface{}, ids ...int64) error {
	if len(ids) == 1 {
		return d.dao.UpdateById(d.tableName, ids[0], fields).Err
	}

	sqb := new(mysql.SqlQueryBuilder)
	sqb.Update(d.tableName).
		WhereConditionAnd(&mysql.SqlColQueryItem{
			Name:      ColumnNameID,
			Condition: mysql.SqlCondIn,
			Value:     ids,
		})

	_, err := d.dao.Exec(sqb.Query(), sqb.Args()...)
	return err
}

func (d *BaseDao) SelectByIDs(dest interface{}, ids ...int64) error {
	if len(ids) == 1 {
		return d.dao.SelectEntityById(d.tableName, ids[0], dest)
	}

	return d.SimpleQueryAnd(&mysql.SqlQueryParams{
		CondItems: []*mysql.SqlColQueryItem{
			{
				Name:      ColumnNameID,
				Condition: mysql.SqlCondIn,
				Value:     ids,
			},
		},
	}, dest)
}

func (d *BaseDao) SimpleQueryAnd(params *mysql.SqlQueryParams, dest interface{}) error {
	return d.dao.SimpleQueryEntitiesAnd(d.tableName, params, dest)
}

func (d *BaseDao) SimpleTotalAnd(items ...*mysql.SqlColQueryItem) (int64, error) {
	return d.dao.SimpleTotalAnd(d.tableName, items...)
}

func (d *BaseDao) Begin() error {
	return d.dao.Begin()
}

func (d *BaseDao) Commit() error {
	return d.dao.Commit()
}

func (d *BaseDao) Rollback() error {
	return d.dao.Rollback()
}
