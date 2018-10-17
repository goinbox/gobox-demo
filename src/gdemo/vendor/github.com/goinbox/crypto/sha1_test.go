package crypto

import "testing"

func TestSha1String(t *testing.T) {
	sha1 := Sha1String([]byte("123"))
	if len(sha1) != 40 {
		t.Error(string(sha1))
	}
	t.Log(string(sha1))
}
