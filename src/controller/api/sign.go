package api

import (
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/goinbox/crypto"
	"github.com/goinbox/gomisc"

	"gdemo/perror"
)

type ApiSignParams struct {
	T     int64
	Nonce string
	Sign  string
}

var ApiSignQueryNames = []string{"t", "nonce"}

// func SetApiSignParams(qs *query.QuerySet, asp *ApiSignParams) {
// qs.Int64Var(&asp.T, "t", true, perror.ECommonInvalidArg, "invalid sign t", query.CheckInt64IsPositive)
// qs.StringVar(&asp.Nonce, "nonce", true, perror.ECommonInvalidArg, "invalid sign nonce", query.CheckStringNotEmpty)
// qs.StringVar(&asp.Sign, "sign", true, perror.ECommonInvalidArg, "invalid sign sign", query.CheckStringNotEmpty)
// qs.IntVar(&asp.Debug, "debug", false, perror.ECommonInvalidArg, "invalid sign debug", nil)
// }

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
