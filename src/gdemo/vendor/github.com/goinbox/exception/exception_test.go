package exception

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	e := New(101, "test exception")

	fmt.Println(e.Errno(), e.Msg())
	fmt.Println(e.Error())
}
