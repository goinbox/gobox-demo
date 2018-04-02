package query

import (
	"github.com/goinbox/exception"
)

type Value interface {
	Required() bool
	Set(str string) error
	Check() bool
	Error() *exception.Exception
}

type baseValue struct {
	required bool
	errno    int
	msg      string
}

func newBaseValue(required bool, errno int, msg string) *baseValue {
	return &baseValue{
		required: required,
		errno:    errno,
		msg:      msg,
	}
}

func (b *baseValue) Required() bool {
	return b.required
}

func (b *baseValue) Error() *exception.Exception {
	return exception.New(b.errno, b.msg)
}
