package mongo

import "time"

type Query struct {
	finder     interface{}
	selector   interface{}
	sort       []string
	limit      int
	skip       int
	setMaxTime time.Duration
}

func (this *Query) Find(finder interface{}) *Query {
	this.finder = finder
	return this
}

func (this *Query) Select(selector interface{}) *Query {
	this.selector = selector
	return this
}

func (this *Query) Sort(fields ...string) *Query {
	this.sort = fields
	return this
}

func (this *Query) Limit(n int) *Query {
	this.limit = n
	return this
}

func (this *Query) Skip(n int) *Query {
	this.skip = n
	return this
}

func (this *Query) SetMaxTime(maxTime time.Duration) *Query {
	this.setMaxTime = maxTime
	return this
}

func NewQuery() *Query {
	query := &Query{}
	return query
}
