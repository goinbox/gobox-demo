package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/goinbox/gohttp/v6/httpserver"
	"github.com/goinbox/router"

	"gdemo/pcontext"
	"gdemo/tracing"

	"gdemo/conf"
	"gdemo/controller/api/demo"
	"gdemo/perror"
	"gdemo/resource"
	"gdemo/validate"
)

var args struct {
	confDir string
}

func main() {
	parseArgs()

	err := conf.Init(args.confDir)
	if err != nil {
		fmt.Println("conf.Init error:", err)
		os.Exit(perror.ESysInitConfError)
	}

	err = resource.InitLog(conf.ServerConf.Log.Api)
	if err != nil {
		fmt.Println("resource.InitLog error:", err)
		os.Exit(perror.ESysInitLogFail)
	}
	defer func() {
		resource.FreeLog()
	}()

	err = resource.InitMySQL(conf.ServerConf.MySQL)
	if err != nil {
		fmt.Println("resource.InitMySQL error:", err)
		os.Exit(perror.ESysMysqlError)
	}

	resource.InitRedis(conf.ServerConf.Redis)

	err = tracing.Init(conf.ServerConf.Tracing, "gdemo-api")
	if err != nil {
		fmt.Println("tracing.Init error:", err)
		os.Exit(perror.ESysInitTracingFail)
	}
	defer func() {
		err = tracing.Shutdown()
		if err != nil {
			fmt.Println("tracing.Shutdown error:", err)
		}
	}()

	pprof(conf.ServerConf.Pprof)
	err = runServer(conf.ServerConf.Api)
	if err != nil {
		fmt.Println("runServer error:", err)
		os.Exit(perror.ESysRunServerError)
	}
}

func parseArgs() {
	flag.StringVar(&args.confDir, "conf-dir", "", "conf-dir absolute path")
	flag.Parse()

	args.confDir = strings.TrimRight(args.confDir, "/")
	if args.confDir == "" {
		fmt.Println("missing flag conf-dir")
		flag.PrintDefaults()
		os.Exit(perror.ESysInvalidConfDir)
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
	resource.AccessLogger.Notice(fmt.Sprintf("runServer %s", addr))

	handler := httpserver.NewHandler[*pcontext.Context](r).
		SetStartTraceFunc(tracing.StartTraceForFramework)
	err = httpserver.NewServer(addr, handler).
		ListenAndServe(pcontext.NewContext(resource.AccessLogger))
	if err != nil {
		return fmt.Errorf("httpserver.ListenAndServe error: %w", err)
	}

	return nil
}
