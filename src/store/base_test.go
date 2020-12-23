package store

import (
	"os"
	"path/filepath"

	"gdemo/conf"
	"gdemo/resource"
)

func init() {
	curDir, _ := os.Getwd()
	prjHome := curDir + "/../../"
	prjHome, _ = filepath.Abs(prjHome)

	_ = conf.Init(prjHome)

	_ = resource.InitLog("test")
	resource.InitRedis()
	resource.InitMysql()
}
