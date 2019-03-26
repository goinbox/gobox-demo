package misc

import (
	"fmt"
	"gdemo/conf"
	"gdemo/resource"
	"github.com/goinbox/color"
	"github.com/goinbox/golog"
	"os"
	"reflect"
	"sync"
)

func PrintComplexObjectForTest(v interface{}) {
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

var initLock sync.Mutex
var testInit bool = false
var TestLogger golog.ILogger = golog.NewSimpleLogger(
	golog.NewConsoleWriter(),
	golog.NewConsoleFormater(golog.NewSimpleFormater())).
	SetLogLevel(golog.LEVEL_DEBUG)

func InitTestSystem() {
	initLock.Lock()
	defer initLock.Unlock()

	if testInit {
		return
	}

	_ = conf.Init(os.Getenv("GOPATH"))
	_ = resource.InitLog("test")
	resource.InitRedis()
	resource.InitMysql()

	testInit = true
}
