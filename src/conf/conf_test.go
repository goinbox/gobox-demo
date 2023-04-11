package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
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
	pretty, _ := json.MarshalIndent(ServerConf, "", "    ")
	t.Log(string(pretty))
}
