package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/goinbox/gohttp/httpserver"
	"github.com/goinbox/gohttp/router"

	"gdemo/controller/api/validate"
	"gdemo/logic/app"
	"gdemo/logic/factory"
	"gdemo/misc"
	"gdemo/test"
)

var server *httpserver.Server
var appItem *app.App

func InitTestServer(cs ...router.Controller) {
	r := router.NewRouter()
	r.MapRouteItems(cs...)

	err := validate.Init()
	if err != nil {
		panic(fmt.Errorf("validate.InitTestServer error: %w", err))
	}

	server = httpserver.NewServer(r)

	appItem = factory.DefaultLogicFactory.AppLogic().ListAllApps(test.Context())[0]
}

func HandleRequest(target, body string) {
	r := httptest.NewRequest(http.MethodPost, target, bytes.NewReader([]byte(body)))
	// r.RemoteAddr = "127.0.0.1:12345"
	ts := time.Now().Unix()
	r.Header.Set("App", appItem.Name)
	r.Header.Set("Timestamp", strconv.FormatInt(ts, 10))
	r.Header.Set("Sign", misc.ApiSign(appItem.Token, ts, []byte(body)))

	w := httptest.NewRecorder()
	server.ServeHTTP(w, r)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("response data error: %w", err))
	}

	fmt.Println(string(data))
}
