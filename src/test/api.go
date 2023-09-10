package test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"gdemo/logic/app"
	"gdemo/misc"
)

type ApiControllerRunner struct {
	Handler http.Handler
	App     *app.App
}

func (r *ApiControllerRunner) Run(target, body string) ([]byte, error) {
	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader([]byte(body)))
	// r.RemoteAddr = "127.0.0.1:12345"
	ts := time.Now().Unix()
	req.Header.Set("App", r.App.Name)
	req.Header.Set("Timestamp", strconv.FormatInt(ts, 10))
	req.Header.Set("Sign", misc.ApiSign(r.App.Token, ts, []byte(body)))

	w := httptest.NewRecorder()
	r.Handler.ServeHTTP(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()

	return io.ReadAll(resp.Body)
}
