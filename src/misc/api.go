package misc

import (
	"gdemo/errno"

	"github.com/goinbox/crypto"
	"github.com/goinbox/goerror"
	"github.com/goinbox/gohttp/query"
	"github.com/goinbox/gomisc"

	"encoding/json"
	"net/url"
	"sort"
	"strings"
	"time"
)

type ApiData struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
	V     string `json:"v"`

	Data interface{} `json:"data"`
}

func ApiJson(v string, data interface{}, e *goerror.Error) []byte {
	result := &ApiData{
		Errno: errno.Success,
		Msg:   "",
		V:     v,

		Data: data,
	}
	if e != nil {
		result.Errno = e.Errno()
		result.Msg = e.Msg()
	}

	aj, err := json.Marshal(result)
	if err != nil {
		result.Errno = errno.ECommonJsonEncodeError
		result.Msg = err.Error()
		result.Data = nil

		aj, _ = json.Marshal(result)
	}

	return aj
}

func ApiJsonp(v string, data interface{}, e *goerror.Error, callback string) []byte {
	return gomisc.AppendBytes(
		[]byte(" "),
		[]byte(callback),
		[]byte("("),
		[]byte(ApiJson(v, data, e)),
		[]byte(");"),
	)
}

type ApiSignParams struct {
	T     int64
	Nonce string
	Sign  string
	Debug int
}

var ApiSignQueryNames = []string{"t", "nonce"}

func SetApiSignParams(qs *query.QuerySet, asp *ApiSignParams) {
	qs.Int64Var(&asp.T, "t", true, errno.ECommonInvalidArg, "invalid sign t", query.CheckInt64IsPositive)
	qs.StringVar(&asp.Nonce, "nonce", true, errno.ECommonInvalidArg, "invalid sign nonce", query.CheckStringNotEmpty)
	qs.StringVar(&asp.Sign, "sign", true, errno.ECommonInvalidArg, "invalid sign sign", query.CheckStringNotEmpty)
	qs.IntVar(&asp.Debug, "debug", false, errno.ECommonInvalidArg, "invalid sign debug", nil)
}

func VerifyApiSign(asp *ApiSignParams, queryValues url.Values, signQueryNames []string, token string) *goerror.Error {
	if time.Now().Unix()-asp.T > 600 {
		return goerror.New(errno.ECommonInvalidArg, "verify sign failed, invalid sign t")
	}

	sign := CalApiSign(queryValues, signQueryNames, token)
	if sign != asp.Sign {
		return goerror.New(errno.ECommonInvalidArg, "verify sign failed, invalid sign sign")
	}

	return nil
}

func CalApiSign(queryValues url.Values, signQueryNames []string, token string) string {
	var sign string
	sort.Strings(signQueryNames)

	for i, name := range signQueryNames {
		value := queryValues.Get(name)
		value = strings.TrimSpace(value)
		if value != "" {
			if i != 0 {
				sign += "&"
			}
			sign += name + "=" + url.QueryEscape(value)
		}
	}

	sign = crypto.Md5String([]byte(sign)) + token
	sign = crypto.Md5String([]byte(sign))
	sign, _ = gomisc.SubString(sign, 3, 7)

	return sign
}
