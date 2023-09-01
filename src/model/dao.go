package model

import (
	"gdemo/pcontext"
	"github.com/goinbox/mysql"
)

type Dao interface {
	Insert(ctx *pcontext.Context, entities ...interface{}) *mysql.SqlExecResult
	DeleteByIDs(ctx *pcontext.Context, ids ...int64) *mysql.SqlExecResult
	UpdateByIDs(ctx *pcontext.Context, updateColumns []*mysql.SqlUpdateColumn, ids ...int64) *mysql.SqlExecResult
	SelectByID(ctx *pcontext.Context, id int64, dest interface{}) error

	SimpleQueryAnd(ctx *pcontext.Context, params *mysql.SqlQueryParams, dest interface{}) error
	SimpleTotalAnd(ctx *pcontext.Context, items ...*mysql.SqlColQueryItem) (int64, error)

	Begin(ctx *pcontext.Context) error
	Commit(ctx *pcontext.Context) error
	Rollback(ctx *pcontext.Context) error
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

func (d *BaseDao) Insert(ctx *pcontext.Context, entities ...interface{}) *mysql.SqlExecResult {
	return d.Dao.InsertEntities(ctx, d.TableName, entities...)
}

func (d *BaseDao) DeleteByIDs(ctx *pcontext.Context, ids ...int64) *mysql.SqlExecResult {
	return d.Dao.DeleteByIDs(ctx, d.TableName, ids...)
}

func (d *BaseDao) UpdateByIDs(ctx *pcontext.Context,
	updateColumns []*mysql.SqlUpdateColumn, ids ...int64) *mysql.SqlExecResult {
	return d.Dao.UpdateByIDs(ctx, d.TableName, updateColumns, ids...)
}

func (d *BaseDao) SelectByID(ctx *pcontext.Context, id int64, dest interface{}) error {
	return d.Dao.SelectEntityByID(ctx, d.TableName, id, dest)
}

func (d *BaseDao) SimpleQueryAnd(ctx *pcontext.Context, params *mysql.SqlQueryParams, dest interface{}) error {
	return d.Dao.SimpleQueryEntitiesAnd(ctx, d.TableName, params, dest)
}

func (d *BaseDao) SimpleTotalAnd(ctx *pcontext.Context, items ...*mysql.SqlColQueryItem) (int64, error) {
	return d.Dao.SimpleTotalAnd(ctx, d.TableName, items...)
}

func (d *BaseDao) Begin(ctx *pcontext.Context) error {
	return d.Dao.Begin(ctx)
}

func (d *BaseDao) Commit(ctx *pcontext.Context) error {
	return d.Dao.Commit(ctx)
}

func (d *BaseDao) Rollback(ctx *pcontext.Context) error {
	return d.Dao.Rollback(ctx)
}
