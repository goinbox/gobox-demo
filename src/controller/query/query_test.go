package query

import (
	"fmt"
	"net/url"
	"testing"
)

func TestParse(t *testing.T) {
	qv, _ := url.ParseQuery("a=1&b=hello&c=64")
	qs := NewQuerySet()

	var a int
	var b string
	var c int64

	qs.IntVar(&a, "a", true, 101, "invalid a", CheckIntIsPositive)
	qs.StringVar(&b, "b", true, 102, "invalid b", CheckStringNotEmpty)
	qs.Int64Var(&c, "c", false, 103, "invalid c", CheckInt64IsPositive)

	e := qs.Parse(qv)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Println(a, b, c)
	}

	fmt.Println(qs.ExistsInfo())
}
