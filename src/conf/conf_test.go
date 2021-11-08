package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/goinbox/color"
)

func init() {
	wd, _ := os.Getwd()
	prjHome := filepath.Dir(filepath.Dir(wd))

	e := Init(prjHome + "/conf/server")
	if e != nil {
		fmt.Println("Init error: ", e.Error())
	}
}

func TestConf(t *testing.T) {
	printComplexObjectForTest(&ServerConf)

	for k, item := range ServerConf.Log {
		t.Log(k, item)
	}

	printComplexObjectForTest(&ServerConf.Pprof)
	printComplexObjectForTest(&ServerConf.Api)
	printComplexObjectForTest(ServerConf.Redis)
	printComplexObjectForTest(ServerConf.MySQL)
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
