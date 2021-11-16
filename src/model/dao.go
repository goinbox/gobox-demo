package model

import "github.com/goinbox/mysql"

type Dao interface {
	Insert(entities ...interface{}) *mysql.SqlExecResult
	DeleteByIDs(ids ...int64) *mysql.SqlExecResult
	UpdateByIDs(fields map[string]interface{}, ids ...int64) *mysql.SqlExecResult
	SelectByID(id int64, dest interface{}) error

	SimpleQueryAnd(params *mysql.SqlQueryParams, dest interface{}) error
	SimpleTotalAnd(items ...*mysql.SqlColQueryItem) (int64, error)

	Begin() error
	Commit() error
	Rollback() error
}

type BaseDao struct {
	TableName string
	Dao       *mysql.EntityDao
}

func NewBaseDao(TableName string, client *mysql.Client) *BaseDao {
	return &BaseDao{
		TableName: TableName,
		Dao: &mysql.EntityDao{
			Dao: mysql.Dao{
				Client: client,
			},
		},
	}
}

func (d *BaseDao) Insert(entities ...interface{}) *mysql.SqlExecResult {
	return d.Dao.InsertEntities(d.TableName, entities...)
}

func (d *BaseDao) DeleteByIDs(ids ...int64) *mysql.SqlExecResult {
	return d.Dao.DeleteByIDs(d.TableName, ids...)
}

func (d *BaseDao) UpdateByIDs(fields map[string]interface{}, ids ...int64) *mysql.SqlExecResult {
	return d.Dao.UpdateByIDs(d.TableName, fields, ids...)
}

func (d *BaseDao) SelectByID(id int64, dest interface{}) error {
	return d.Dao.SelectEntityByID(d.TableName, id, dest)
}

func (d *BaseDao) SimpleQueryAnd(params *mysql.SqlQueryParams, dest interface{}) error {
	return d.Dao.SimpleQueryEntitiesAnd(d.TableName, params, dest)
}

func (d *BaseDao) SimpleTotalAnd(items ...*mysql.SqlColQueryItem) (int64, error) {
	return d.Dao.SimpleTotalAnd(d.TableName, items...)
}

func (d *BaseDao) Begin() error {
	return d.Dao.Begin()
}

func (d *BaseDao) Commit() error {
	return d.Dao.Commit()
}

func (d *BaseDao) Rollback() error {
	return d.Dao.Rollback()
}
