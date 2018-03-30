/**
* @file exception.go
* @brief exception with errno and msg
* @author ligang
* @date 2016-02-01
 */

package exception

import (
	"strconv"
)

type Exception struct {
	errno int
	msg   string
}

func New(errno int, msg string) *Exception {
	return &Exception{
		errno: errno,
		msg:   msg,
	}
}

func (e *Exception) Error() string {
	result := "errno: " + strconv.Itoa(e.errno) + ", "
	result += "msg: " + e.msg

	return result
}

func (e *Exception) Errno() int {
	return e.errno
}

func (e *Exception) Msg() string {
	return e.msg
}
