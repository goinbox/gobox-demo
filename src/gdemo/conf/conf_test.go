package conf

import (
	"gdemo/misc"

	"fmt"
	"os"
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
	misc.PrintComplexObjectForTest(&BaseConf)
	misc.PrintComplexObjectForTest(&LogConf)
	misc.PrintComplexObjectForTest(&PprofConf)
	misc.PrintComplexObjectForTest(&ApiHttpConf)

	for _, item := range RedisConfList {
		misc.PrintComplexObjectForTest(item)
	}

	misc.PrintComplexObjectForTest(&MysqlConf)
	misc.PrintComplexObjectForTest(&MongoConf)
}
