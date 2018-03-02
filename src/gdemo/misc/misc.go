package misc

import (
	"gdemo/errno"

	"github.com/goinbox/exception"
	"github.com/goinbox/gomisc"

	"encoding/json"
	"net/smtp"
	"reflect"
	"strings"
)

type ApiData struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
	V     string `json:"v"`

	Data interface{} `json:"data"`
}

func ApiJson(v string, data interface{}, e *exception.Exception) []byte {
	result := &ApiData{
		Errno: errno.SUCCESS,
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
		result.Errno = errno.E_COMMON_JSON_ENCODE_ERROR
		result.Msg = err.Error()
		result.Data = nil

		aj, _ = json.Marshal(result)
	}

	return aj
}

func ApiJsonp(v string, data interface{}, e *exception.Exception, callback string) []byte {
	return gomisc.AppendBytes(
		[]byte(" "),
		[]byte(callback),
		[]byte("("),
		[]byte(ApiJson(v, data, e)),
		[]byte(");"),
	)
}

func SendMail(subject, body, from string, to []string) error {
	auth := smtp.PlainAuth("", "", "", "")
	toStr := strings.Join(to, "\t")
	msg := []byte("To:" + toStr + "\r\n" + "Subject:" + subject + "\r\n\r\n" + body + "\r\n")

	err := smtp.SendMail("127.0.0.1:25", auth, from, to, msg)
	return err
}

var StructSimpleFields map[reflect.Kind]bool = map[reflect.Kind]bool{
	reflect.Bool:    true,
	reflect.Int:     true,
	reflect.Int8:    true,
	reflect.Int16:   true,
	reflect.Int32:   true,
	reflect.Int64:   true,
	reflect.Uint:    true,
	reflect.Uint8:   true,
	reflect.Uint16:  true,
	reflect.Uint64:  true,
	reflect.Float32: true,
	reflect.Float64: true,
	reflect.String:  true,
}

func StructSimpleFieldAssign(sou, dst interface{}) {
	rsoue := reflect.ValueOf(sou).Elem()
	rsout := rsoue.Type()
	rdste := reflect.ValueOf(dst).Elem()

	for i := 0; i < rsoue.NumField(); i++ {
		rsoutf := rsout.Field(i)
		rdstv := rdste.FieldByName(rsoutf.Name)
		if !rdstv.IsValid() {
			continue
		}

		rsouv := rsoue.Field(i)
		rsouk := rsouv.Kind()
		if _, ok := StructSimpleFields[rsouk]; ok && rsouv.Type().Name() == rdstv.Type().Name() {
			rdstv.Set(rsouv)
		}
	}
}
