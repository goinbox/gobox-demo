package conf

import (
	"github.com/goinbox/color"

	"fmt"
	"os"
	"reflect"
	"testing"
)

func init() {
	prjHome := os.Getenv("GOPATH")

	e := Init(prjHome)
	if e != nil {
		fmt.Println("Init error: ", e.Error())
	}
}

func TestConf(t *testing.T) {
	t.Log("PrjHome", PrjHome)
	printComplexObjectForTest(&BaseConf)
	printComplexObjectForTest(&LogConf)
	printComplexObjectForTest(&PprofConf)
	printComplexObjectForTest(&ApiHttpConf)

	for _, item := range RedisConfList {
		printComplexObjectForTest(item)
	}

	printComplexObjectForTest(&MysqlConf)
	printComplexObjectForTest(&MongoConf)
}

func printComplexObjectForTest(v interface{}) {
	vo := reflect.ValueOf(v)
	elems := vo.Elem()
	ts := elems.Type()

	c := color.Yellow([]byte("Print detail: "))
	fmt.Println(string(c), vo.Type())
	for i := 0; i < elems.NumField(); i++ {
		field := elems.Field(i)
		fmt.Println(ts.Field(i).Name, field.Type(), field.Interface())
	}
}
