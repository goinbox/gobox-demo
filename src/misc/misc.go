package misc

import (
	"crypto/sha256"
	"encoding/hex"
	"net/smtp"
	"reflect"
	"strings"
)

func SendMail(subject, body, from string, to []string) error {
	auth := smtp.PlainAuth("", "", "", "")
	toStr := strings.Join(to, "\t")
	msg := []byte("To:" + toStr + "\r\n" + "Subject:" + subject + "\r\n\r\n" + body + "\r\n")

	err := smtp.SendMail("127.0.0.1:25", auth, from, to, msg)
	return err
}

var structSimpleFields = map[reflect.Kind]bool{
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
		if _, ok := structSimpleFields[rsouk]; ok && rsouv.Type().Name() == rdstv.Type().Name() {
			rdstv.Set(rsouv)
		}
	}
}

func Sha256(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}
