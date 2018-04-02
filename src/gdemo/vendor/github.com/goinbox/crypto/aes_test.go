package crypto

import (
	"testing"
)

func TestAesCBCCrypter(t *testing.T) {
	key := Md5String([]byte("gobox"))
	iv := Md5String([]byte("andals"))[:AES_BLOCK_SIZE]
	data := []byte("abc")

	acc, err := NewAesCBCCrypter([]byte(key), []byte(iv))
	t.Log(err)
	t.Log(acc.BlockSize())

	crypted := acc.Encrypt(data)
	t.Log(crypted)

	d := acc.Decrypt(crypted)
	t.Log(d)

	if string(d) != string(data) {
		t.Error(d, data)
	}
}
