package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/goinbox/gomisc"
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
	pretty, _ := gomisc.PrettyJson(ServerConf)
	t.Log(string(pretty))
}
