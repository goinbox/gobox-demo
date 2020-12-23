package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"

	"github.com/goinbox/gohttp/gracehttp"
	"github.com/goinbox/gohttp/router"
	"github.com/goinbox/gohttp/system"
	"github.com/goinbox/pidfile"

	"gdemo/conf"
	"gdemo/controller/api/demo"
	"gdemo/controller/api/mongodemo"
	"gdemo/errno"
	"gdemo/resource"
)

func main() {
	var prjHome string

	flag.StringVar(&prjHome, "prj-home", "", "prj-home absolute path")
	flag.Parse()

	prjHome = strings.TrimRight(prjHome, "/")
	if prjHome == "" {
		fmt.Println("missing flag prjHome: ")
		flag.PrintDefaults()
		os.Exit(errno.ESysInvalidPrjHome)
	}

	e := conf.Init(prjHome)
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(e.Errno())
	}

	if conf.PprofConf.Enable {
		go func() {
			_ = http.ListenAndServe("127.0.0.1:"+conf.PprofConf.Port, nil)
		}()
	}

	e = resource.InitLog("api")
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(e.Errno())
	}
	defer func() {
		resource.FreeLog()
	}()

	resource.InitMysql()
	resource.InitRedis()
	resource.InitMongo()

	pf, err := pidfile.CreatePidFile(conf.BaseConf.ApiPidFile)
	if err != nil {
		fmt.Printf("create pid file %s failed, error: %s\n", conf.BaseConf.ApiPidFile, err.Error())
		os.Exit(errno.ESysSavePidFileFail)
	}

	r := router.NewSimpleRouter()
	r.MapRouteItems(
		new(demo.DemoController),
		new(mongodemo.MongoDemoController),
	)

	sys := system.NewSystem(r)

	err = gracehttp.ListenAndServe(conf.ApiHttpConf.GoHttpHost+":"+conf.ApiHttpConf.GoHttpPort, sys)
	if err != nil {
		fmt.Println("pid:" + strconv.Itoa(os.Getpid()) + ", err:" + err.Error())
	}

	if err := pidfile.ClearPidFile(pf); err != nil {
		fmt.Printf("clear pid file failed, error: %s\n", err.Error())
	}
}
