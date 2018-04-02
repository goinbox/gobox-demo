package query

import (
	"github.com/goinbox/exception"

	"net/url"
	"strings"
)

type QuerySet struct {
	formal map[string]Value
	exists map[string]bool
}

func NewQuerySet() *QuerySet {
	q := &QuerySet{
		formal: make(map[string]Value),
		exists: make(map[string]bool),
	}

	return q
}

func (q *QuerySet) ExistsInfo() map[string]bool {
	return q.exists
}

func (q *QuerySet) Exist(name string) bool {
	if v, ok := q.exists[name]; ok && v {
		return true
	}

	return false
}

func (q *QuerySet) Var(name string, v Value) *QuerySet {
	q.formal[name] = v

	return q
}

func (q *QuerySet) IntVar(p *int, name string, required bool, errno int, msg string, cf CheckInt) *QuerySet {
	q.Var(name, NewIntValue(p, required, errno, msg, cf))

	return q
}

func (q *QuerySet) StringVar(p *string, name string, required bool, errno int, msg string, cf CheckString) *QuerySet {
	q.Var(name, NewStringValue(p, required, errno, msg, cf))

	return q
}

func (q *QuerySet) Int64Var(p *int64, name string, required bool, errno int, msg string, cf CheckInt64) *QuerySet {
	q.Var(name, NewInt64Value(p, required, errno, msg, cf))

	return q
}

func (q *QuerySet) Parse(actual url.Values) *exception.Exception {
	for name, v := range q.formal {
		if len(actual[name]) == 0 {
			if v.Required() {
				return v.Error()
			}
			continue
		}

		q.exists[name] = true
		str := strings.TrimSpace(actual.Get(name))
		err := v.Set(str)
		if err != nil {
			return v.Error()
		}
		if v.Check() == false {
			return v.Error()
		}
	}

	return nil
}
