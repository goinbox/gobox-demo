package crypto

import (
	"crypto/md5"
	"fmt"
)

func Md5String(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}
