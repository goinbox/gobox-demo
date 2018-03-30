package client

import (
	"github.com/goinbox/golog"

	"fmt"
	"testing"
	"time"
)

func TestClientGet(t *testing.T) {
	path := "/tmp/test_http_client.log"
	w, _ := golog.NewFileWriter(path)
	logger, _ := golog.NewSimpleLogger(w, golog.LEVEL_INFO, nil)

	client := NewClient(logger).
		SetTimeout(time.Second * 3).
		SetMaxIdleConnsPerHost(10)
	extHeaders := map[string]string{
		"GO-CLIENT-1": "gobox-httpclient-1",
		"GO-CLIENT-2": "gobox-httpclient-2",
	}
	req, _ := NewRequestForGet("http://www.vmc7.com/test.php", "127.0.0.1", extHeaders)

	resp, err := client.Do(req, 1)
	fmt.Println(string(resp.Contents), resp.T.String(), err)

	req, _ = NewRequestForGet("http://www.vmc7.com/index.html", "127.0.0.1", extHeaders)

	resp, err = client.Do(req, 1)
	fmt.Println(string(resp.Contents), resp.T.String(), err)
}
