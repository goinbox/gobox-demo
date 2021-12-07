package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/goinbox/gohttp/gracehttp"
	"github.com/goinbox/gohttp/httpserver"
	"github.com/goinbox/gohttp/router"

	"gdemo/conf"
	"gdemo/controller/api/demo"
	"gdemo/controller/api/validate"
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

	err = resource.InitLog(conf.ServerConf.Log.Api)
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
	err = runServer(conf.ServerConf.Api)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(perror.ESysRunServerError)
	}
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

func runServer(config *conf.ApiConf) error {
	r := router.NewRouter()
	r.MapRouteItems(
		new(demo.Controller),
	)

	err := validate.Init()
	if err != nil {
		return fmt.Errorf("validate.Init error: %w", err)
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	err = gracehttp.ListenAndServe(addr, httpserver.NewServer(r))
	if err != nil {
		return fmt.Errorf("ListenAndServe error: %w", err)
	}

	return nil
}
