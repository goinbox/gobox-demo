package misc

import (
	"github.com/goinbox/color"
	"github.com/goinbox/mysql"
	"github.com/goinbox/redis"
	"github.com/goinbox/golog"

	"fmt"
	"reflect"
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

func MysqlTestClient() *mysql.Client {
	config := mysql.NewConfig("root", "123", "127.0.0.1", "3306", "gobox-demo")

	w, _ := golog.NewFileWriter("/tmp/test_mysql.log")
	logger, _ := golog.NewSimpleLogger(w, golog.LEVEL_INFO, golog.NewSimpleFormater())

	client, _ := mysql.NewClient(config, logger)

	return client
}

func RedisTestClient() *redis.Client {
	config := redis.NewConfig("127.0.0.1", "6379", "123")

	w, _ := golog.NewFileWriter("/tmp/test_redis.log")
	logger, _ := golog.NewSimpleLogger(w, golog.LEVEL_INFO, golog.NewSimpleFormater())

	return redis.NewClient(config, logger)
}
