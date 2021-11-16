package query

import "strconv"

type CheckInt func(v int) bool

func CheckIntIsPositive(v int) bool {
	if v > 0 {
		return true
	}
	return false
}

type intValue struct {
	*baseValue

	p  *int
	cf CheckInt
}

func NewIntValue(p *int, required bool, errno int, msg string, cf CheckInt) *intValue {
	i := &intValue{
		baseValue: newBaseValue(required, errno, msg),

		p:  p,
		cf: cf,
	}

	return i
}

func (i *intValue) Set(str string) error {
	var v int = 0
	var e error = nil

	if str != "" {
		v, e = strconv.Atoi(str)
	}

	if e != nil {
		return e
	}

	*(i.p) = v

	return nil
}

func (i *intValue) Check() bool {
	if i.cf == nil {
		return true
	}

	return i.cf(*(i.p))
}
