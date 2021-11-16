package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/goinbox/gohttp/gracehttp"
	"github.com/goinbox/gohttp/router"
	"github.com/goinbox/gohttp/system"
	"github.com/goinbox/pidfile"

	"gdemo/conf"
	"gdemo/controller/api/demo"
	"gdemo/perror"
	"gdemo/resource"
)

var args struct {
	prjHome string
}

func main() {
	parseArgs()

	err := conf.Init(args.prjHome + "/conf/server")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(perror.ESysInitConfError)
	}

	err = resource.InitLog(conf.ServerConf.Log["api"])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(perror.ESysInitLogFail)
	}
	defer func() {
		resource.FreeLog()
	}()

	err = resource.InitMySQL(conf.ServerConf.MySQL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(perror.ESysMysqlError)
	}

	resource.InitRedis(conf.ServerConf.Redis)

	pprof(conf.ServerConf.Pprof)
	runServer(conf.ServerConf.Api)
}

func parseArgs() {
	flag.StringVar(&args.prjHome, "prj-home", "", "prj-home absolute path")
	flag.Parse()

	args.prjHome = strings.TrimRight(args.prjHome, "/")
	if args.prjHome == "" {
		fmt.Println("missing flag prj-home: ")
		flag.PrintDefaults()
		os.Exit(perror.ESysInvalidPrjHome)
	}
}

func pprof(config *conf.PprofConf) {
	if config.Enable {
		go func() {
			addr := fmt.Sprintf("127.0.0.1:%d", config.Port)
			_ = http.ListenAndServe(addr, nil)
		}()
	}
}

func runServer(config *conf.ApiConf) {
	pf, err := pidfile.CreatePidFile(config.PidFile)
	if err != nil {
		fmt.Printf("create pid file %s failed, error: %v\n", config.PidFile, err)
		os.Exit(perror.ESysFileIOError)
	}

	r := router.NewSimpleRouter()
	r.MapRouteItems(
		new(demo.DemoController),
	)

	sys := system.NewSystem(r)

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	err = gracehttp.ListenAndServe(addr, sys)
	if err != nil {
		fmt.Printf("pid:%d error: %v", os.Getpid(), err)
	}

	if err := pidfile.ClearPidFile(pf); err != nil {
		fmt.Printf("clear pid file failed, error: %v\n", err)
	}
}
