package misc

import (
	"net/url"
	"testing"
)

func TestCalApiSign(t *testing.T) {
	qvs, _ := url.ParseQuery("pid=132&startTime=1519747200&endTime=1519920000&t=123&nonce=abc&token=ghi")
	qns := []string{"pid", "startTime", "endTime", "t", "nonce"}
	token := "ghi"

	sign := CalApiSign(qvs, qns, token)
	t.Log(sign)
}
