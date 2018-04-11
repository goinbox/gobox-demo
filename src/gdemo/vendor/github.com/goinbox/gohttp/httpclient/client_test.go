package httpclient

import (
	"github.com/goinbox/golog"

	"net/http"
	"testing"
	"time"
)

func TestClientGet(t *testing.T) {
	client := getClient()
	extHeaders := map[string]string{
		"GO-CLIENT-1": "gobox-httpclient-1",
		"GO-CLIENT-2": "gobox-httpclient-2",
	}
	req, _ := NewRequest(http.MethodGet, "http://www.vmu16.com/test.php?a=1&b=2", nil, "127.0.0.1", extHeaders)

	resp, err := client.Do(req, 1)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(resp.Contents), resp.T.String())
	}
}

func TestClientPost(t *testing.T) {
	client := getClient()
	extHeaders := map[string]string{
		"GO-CLIENT-1":  "gobox-httpclient-1",
		"GO-CLIENT-2":  "gobox-httpclient-2",
		"Content-Type": "application/x-www-form-urlencoded;charset=utf-8",
	}
	params := map[string]interface{}{
		"a": 1,
		"b": "bb",
		"c": "测试post",
	}
	req, _ := NewRequest(http.MethodPost, "http://www.vmu16.com/test.php", MakeRequestBodyUrlEncoded(params), "127.0.0.1", extHeaders)

	resp, err := client.Do(req, 1)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(resp.Contents), resp.T.String())
	}
}

func getClient() *Client {
	path := "/tmp/test_http_client.log"
	w, _ := golog.NewFileWriter(path)
	logger, _ := golog.NewSimpleLogger(w, golog.LEVEL_INFO, golog.NewSimpleFormater())

	config := NewConfig()
	config.Timeout = time.Second * 1

	return NewClient(config, logger)
}
