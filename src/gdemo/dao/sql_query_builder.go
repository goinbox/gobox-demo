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

func (s *SqlQueryBuilder) Query() string {
	return s.query
}

func (s *SqlQueryBuilder) Args() []interface{} {
	return s.args
}

func (s *SqlQueryBuilder) Insert(tableName string, colNames ...string) *SqlQueryBuilder {
	s.args = nil

	s.query = "INSERT INTO " + tableName + " ("
	s.query += strings.Join(colNames, ", ") + ")"

	return s
}

func (s *SqlQueryBuilder) Values(colsValues ...[]interface{}) *SqlQueryBuilder {
	l := len(colsValues) - 1
	if l == -1 {
		return s
	}

	s.query += " VALUES "
	for i := 0; i < l; i++ {
		s.buildColValues(colsValues[i])
		s.query += ", "
	}
	s.buildColValues(colsValues[l])

	return s
}

func (s *SqlQueryBuilder) Delete(tableName string) *SqlQueryBuilder {
	s.args = nil

	s.query = "DELETE FROM " + tableName

	return s
}

func (s *SqlQueryBuilder) Update(tableName string) *SqlQueryBuilder {
	s.args = nil

	s.query = "UPDATE " + tableName

	return s
}

func (s *SqlQueryBuilder) Set(setItems ...*SqlColQueryItem) *SqlQueryBuilder {
	l := len(setItems) - 1
	if l == -1 {
		return s
	}

	s.query += " SET "
	for i := 0; i < l; i++ {
		s.query += setItems[i].Name + " = ?, "
		s.args = append(s.args, setItems[i].Value)
	}
	s.query += setItems[l].Name + " = ? "
	s.args = append(s.args, setItems[l].Value)

	return s
}

func (s *SqlQueryBuilder) Select(what, tableName string) *SqlQueryBuilder {
	s.args = nil

	s.query = "SELECT " + what + " FROM " + tableName

	return s
}

func (s *SqlQueryBuilder) WhereConditionAnd(condItems ...*SqlColQueryItem) *SqlQueryBuilder {
	if len(condItems) == 0 {
		return s
	}

	s.query += " WHERE "

	s.buildWhereCondition("AND", condItems...)

	return s
}

func (s *SqlQueryBuilder) WhereConditionOr(condItems ...*SqlColQueryItem) *SqlQueryBuilder {
	if len(condItems) == 0 {
		return s
	}

	s.query += " WHERE "

	s.buildWhereCondition("OR", condItems...)

	return s
}

func (s *SqlQueryBuilder) OrderBy(orderBy string) *SqlQueryBuilder {
	if orderBy != "" {
		s.query += " ORDER BY " + orderBy
	}

	return s
}

func (s *SqlQueryBuilder) GroupBy(groupBy string) *SqlQueryBuilder {
	if groupBy != "" {
		s.query += " GROUP BY " + groupBy
	}

	return s
}

func (s *SqlQueryBuilder) HavingConditionAnd(condItems ...*SqlColQueryItem) *SqlQueryBuilder {
	if len(condItems) == 0 {
		return s
	}

	s.query += " HAVING "

	s.buildWhereCondition("AND", condItems...)

	return s
}

func (s *SqlQueryBuilder) HavingConditionOr(condItems ...*SqlColQueryItem) *SqlQueryBuilder {
	if len(condItems) == 0 {
		return s
	}

	s.query += " HAVING "

	s.buildWhereCondition("OR", condItems...)

	return s
}

func (s *SqlQueryBuilder) Limit(offset, cnt int64) *SqlQueryBuilder {
	if cnt <= 0 || offset < 0 {
		return s
	}

	s.query += " LIMIT ?, ?"
	s.args = append(s.args, offset, cnt)

	return s
}

func (s *SqlQueryBuilder) buildColValues(colValues []interface{}) {
	l := len(colValues) - 1
	if l == -1 {
		return
	}

	s.query += "("

	for i := 0; i < l; i++ {
		s.query += "?, "
		s.args = append(s.args, colValues[i])
	}

	s.query += "?)"
	s.args = append(s.args, colValues[l])
}

func (s *SqlQueryBuilder) buildWhereCondition(andOr string, condItems ...*SqlColQueryItem) {
	l := len(condItems) - 1
	if l == -1 {
		return
	}

	for i := 0; i < l; i++ {
		s.buildCondition(condItems[i])
		s.query += " " + andOr + " "
	}
	s.buildCondition(condItems[l])
}

func (s *SqlQueryBuilder) buildCondition(condItem *SqlColQueryItem) {
	switch condItem.Condition {
	case SQL_COND_EQUAL, SQL_COND_NOT_EQUAL, SQL_COND_LESS, SQL_COND_LESS_EQUAL, SQL_COND_GREATER, SQL_COND_GREATER_EQUAL:
		s.query += condItem.Name + " " + condItem.Condition + " ?"
		s.args = append(s.args, condItem.Value)
	case SQL_COND_IN:
		s.buildConditionInOrNotIn(condItem, "IN")
	case SQL_COND_NOT_IN:
		s.buildConditionInOrNotIn(condItem, "NOT IN")
	case SQL_COND_LIKE:
		s.query += condItem.Name + " LIKE ?"
		s.args = append(s.args, condItem.Value)
	case SQL_COND_BETWEEN:
		rev := reflect.ValueOf(condItem.Value)
		s.query += condItem.Name + " BETWEEN ? AND ?"
		s.args = append(s.args, rev.Index(0).Interface(), rev.Index(1).Interface())
	}
}

func (s *SqlQueryBuilder) buildConditionInOrNotIn(condItem *SqlColQueryItem, inOrNotIn string) {
	rev := reflect.ValueOf(condItem.Value)
	l := rev.Len() - 1
	if l == -1 {
		return
	}

	s.query += condItem.Name + " " + inOrNotIn + " ("
	for i := 0; i < l; i++ {
		s.query += "?, "
	}
	s.query += "?)"

	for i := 0; i < rev.Len(); i++ {
		s.args = append(s.args, rev.Index(i).Interface())
	}
}
