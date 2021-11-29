package api

import (
	"gdemo/controller/query"
	"gdemo/perror"

	"github.com/goinbox/crypto"
	"github.com/goinbox/gomisc"

	"encoding/json"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Response struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
	Tid   string `json:"tid"`

	Data interface{} `json:"data"`
}

func ApiJson(data *ApiData) []byte {
	result := &Response{
		Errno: perror.Success,
		Msg:   "",
		Tid:   data.Tid,

		Data: data.Data,
	}
	if data.Err != nil {
		result.Errno = data.Err.Errno()
		result.Msg = data.Err.Msg()
	}

	aj, _ := json.Marshal(result)

	return aj
}

func ApiJsonp(data *ApiData, callback string) []byte {
	return gomisc.AppendBytes(
		[]byte(" "),
		[]byte(callback),
		[]byte("("),
		ApiJson(data),
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
	qs.Int64Var(&asp.T, "t", true, perror.ECommonInvalidArg, "invalid sign t", query.CheckInt64IsPositive)
	qs.StringVar(&asp.Nonce, "nonce", true, perror.ECommonInvalidArg, "invalid sign nonce", query.CheckStringNotEmpty)
	qs.StringVar(&asp.Sign, "sign", true, perror.ECommonInvalidArg, "invalid sign sign", query.CheckStringNotEmpty)
	qs.IntVar(&asp.Debug, "debug", false, perror.ECommonInvalidArg, "invalid sign debug", nil)
}

func VerifyApiSign(asp *ApiSignParams, queryValues url.Values, signQueryNames []string, token string) *perror.Error {
	if time.Now().Unix()-asp.T > 600 {
		return perror.New(perror.ECommonInvalidArg, "verify sign failed, invalid sign t")
	}

	sign := CalApiSign(queryValues, signQueryNames, token)
	if sign != asp.Sign {
		return perror.New(perror.ECommonInvalidArg, "verify sign failed, invalid sign sign")
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
