package misc

import (
	"testing"
	"time"
)

func TestStructSimpleFieldAssign(t *testing.T) {
	type ta struct {
		Id   int64
		Name string
	}
	type tb struct {
		Id   time.Duration
		Name string
	}

	a := &ta{
		Id:   10,
		Name: "a",
	}
	b := new(tb)

	StructSimpleFieldAssign(a, b)
	t.Log(a, b)
}
