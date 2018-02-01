package dao

import (
	"reflect"
	"strings"
)

const (
	SQL_COND_EQUAL         = "="
	SQL_COND_NOT_EQUAL     = "!="
	SQL_COND_LESS          = "<"
	SQL_COND_LESS_EQUAL    = "<="
	SQL_COND_GREATER       = ">"
	SQL_COND_GREATER_EQUAL = ">="
	SQL_COND_IN            = "in"
	SQL_COND_NOT_IN        = "not in"
	SQL_COND_LIKE          = "like"
	SQL_COND_BETWEEN       = "between"
)

type SqlColQueryItem struct {
	Name      string
	Condition string
	Value     interface{}
}

func NewSqlColQueryItem(name, condition string, value interface{}) *SqlColQueryItem {
	return &SqlColQueryItem{
		Name:      name,
		Condition: condition,
		Value:     value,
	}
}

type SqlQueryBuilder struct {
	query string
	args  []interface{}
}

func (this *SqlQueryBuilder) Query() string {
	return this.query
}

func (this *SqlQueryBuilder) Args() []interface{} {
	return this.args
}

func (this *SqlQueryBuilder) Insert(tableName string, colNames ...string) *SqlQueryBuilder {
	this.args = nil

	this.query = "INSERT INTO " + tableName + " ("
	this.query += strings.Join(colNames, ", ") + ")"

	return this
}

func (this *SqlQueryBuilder) Values(colsValues ...[]interface{}) *SqlQueryBuilder {
	l := len(colsValues) - 1
	if l == -1 {
		return this
	}

	this.query += " VALUES "
	for i := 0; i < l; i++ {
		this.buildColValues(colsValues[i])
		this.query += ", "
	}
	this.buildColValues(colsValues[l])

	return this
}

func (this *SqlQueryBuilder) Delete(tableName string) *SqlQueryBuilder {
	this.args = nil

	this.query = "DELETE FROM " + tableName

	return this
}

func (this *SqlQueryBuilder) Update(tableName string) *SqlQueryBuilder {
	this.args = nil

	this.query = "UPDATE " + tableName

	return this
}

func (this *SqlQueryBuilder) Set(setItems ...*SqlColQueryItem) *SqlQueryBuilder {
	l := len(setItems) - 1
	if l == -1 {
		return this
	}

	this.query += " SET "
	for i := 0; i < l; i++ {
		this.query += setItems[i].Name + " = ?, "
		this.args = append(this.args, setItems[i].Value)
	}
	this.query += setItems[l].Name + " = ? "
	this.args = append(this.args, setItems[l].Value)

	return this
}

func (this *SqlQueryBuilder) Select(what, tableName string) *SqlQueryBuilder {
	this.args = nil

	this.query = "SELECT " + what + " FROM " + tableName

	return this
}

func (this *SqlQueryBuilder) WhereConditionAnd(condItems ...*SqlColQueryItem) *SqlQueryBuilder {
	if len(condItems) == 0 {
		return this
	}

	this.query += " WHERE "

	this.buildWhereCondition("AND", condItems...)

	return this
}

func (this *SqlQueryBuilder) WhereConditionOr(condItems ...*SqlColQueryItem) *SqlQueryBuilder {
	if len(condItems) == 0 {
		return this
	}

	this.query += " WHERE "

	this.buildWhereCondition("OR", condItems...)

	return this
}

func (this *SqlQueryBuilder) OrderBy(orderBy string) *SqlQueryBuilder {
	if orderBy != "" {
		this.query += " ORDER BY " + orderBy
	}

	return this
}

func (this *SqlQueryBuilder) GroupBy(groupBy string) *SqlQueryBuilder {
	if groupBy != "" {
		this.query += " GROUP BY " + groupBy
	}

	return this
}

func (this *SqlQueryBuilder) HavingConditionAnd(condItems ...*SqlColQueryItem) *SqlQueryBuilder {
	if len(condItems) == 0 {
		return this
	}

	this.query += " HAVING "

	this.buildWhereCondition("AND", condItems...)

	return this
}

func (this *SqlQueryBuilder) HavingConditionOr(condItems ...*SqlColQueryItem) *SqlQueryBuilder {
	if len(condItems) == 0 {
		return this
	}

	this.query += " HAVING "

	this.buildWhereCondition("OR", condItems...)

	return this
}

func (this *SqlQueryBuilder) Limit(offset, cnt int64) *SqlQueryBuilder {
	if cnt <= 0 || offset < 0 {
		return this
	}

	this.query += " LIMIT ?, ?"
	this.args = append(this.args, offset, cnt)

	return this
}

func (this *SqlQueryBuilder) buildColValues(colValues []interface{}) {
	l := len(colValues) - 1
	if l == -1 {
		return
	}

	this.query += "("

	for i := 0; i < l; i++ {
		this.query += "?, "
		this.args = append(this.args, colValues[i])
	}

	this.query += "?)"
	this.args = append(this.args, colValues[l])
}

func (this *SqlQueryBuilder) buildWhereCondition(andOr string, condItems ...*SqlColQueryItem) {
	l := len(condItems) - 1
	if l == -1 {
		return
	}

	for i := 0; i < l; i++ {
		this.buildCondition(condItems[i])
		this.query += " " + andOr + " "
	}
	this.buildCondition(condItems[l])
}

func (this *SqlQueryBuilder) buildCondition(condItem *SqlColQueryItem) {
	switch condItem.Condition {
	case SQL_COND_EQUAL, SQL_COND_NOT_EQUAL, SQL_COND_LESS, SQL_COND_LESS_EQUAL, SQL_COND_GREATER, SQL_COND_GREATER_EQUAL:
		this.query += condItem.Name + " " + condItem.Condition + " ?"
		this.args = append(this.args, condItem.Value)
	case SQL_COND_IN:
		this.buildConditionInOrNotIn(condItem, "IN")
	case SQL_COND_NOT_IN:
		this.buildConditionInOrNotIn(condItem, "NOT IN")
	case SQL_COND_LIKE:
		this.query += condItem.Name + " LIKE ?"
		this.args = append(this.args, condItem.Value)
	case SQL_COND_BETWEEN:
		rev := reflect.ValueOf(condItem.Value)
		this.query += condItem.Name + " BETWEEN ? AND ?"
		this.args = append(this.args, rev.Index(0).Interface(), rev.Index(1).Interface())
	}
}

func (this *SqlQueryBuilder) buildConditionInOrNotIn(condItem *SqlColQueryItem, inOrNotIn string) {
	rev := reflect.ValueOf(condItem.Value)
	l := rev.Len() - 1
	if l == -1 {
		return
	}

	this.query += condItem.Name + " " + inOrNotIn + " ("
	for i := 0; i < l; i++ {
		this.query += "?, "
	}
	this.query += "?)"

	for i := 0; i < rev.Len(); i++ {
		this.args = append(this.args, rev.Index(i).Interface())
	}
}
