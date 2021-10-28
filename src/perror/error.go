package perror

import (
	"fmt"
)

type Error struct {
	errno int
	msg   string
}

func New(errno int, msg string) *Error {
	return &Error{
		errno: errno,
		msg:   msg,
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
