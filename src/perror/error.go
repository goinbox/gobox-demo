package perror

import (
	"errors"
	"fmt"
)

type Error struct {
	errno int
	msg   string
}

func New(errno int, format string, args ...interface{}) *Error {
	return &Error{
		errno: errno,
		msg:   fmt.Sprintf(format, args...),
	}
}

func NewFromError(errno int, err error) *Error {
	return &Error{
		errno: errno,
		msg:   err.Error(),
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("errno: %d, msg: %s", e.errno, e.msg)
}

func (e *Error) Errno() int {
	return e.errno
}

func (e *Error) Msg() string {
	return e.msg
}

func ParsePerror(err error) *Error {
	perr := &Error{}
	if errors.As(err, &perr) {
		return perr
	}
	return nil
}
