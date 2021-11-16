package query

import "strconv"

type CheckInt64 func(v int64) bool

func CheckInt64IsPositive(v int64) bool {
	if v > 0 {
		return true
	}
	return false
}

type int64Value struct {
	*baseValue

	p  *int64
	cf CheckInt64
}

func NewInt64Value(p *int64, required bool, errno int, msg string, cf CheckInt64) *int64Value {
	i := &int64Value{
		baseValue: newBaseValue(required, errno, msg),

		p:  p,
		cf: cf,
	}

	return i
}

func (i *int64Value) Set(str string) error {
	var v int64 = 0
	var e error = nil

	if str != "" {
		v, e = strconv.ParseInt(str, 10, 64)
	}

	if e != nil {
		return e
	}

	*(i.p) = v

	return nil
}

func (i *int64Value) Check() bool {
	if i.cf == nil {
		return true
	}

	return i.cf(*(i.p))
}
