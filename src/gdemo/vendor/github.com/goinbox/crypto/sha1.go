package crypto

import (
	"crypto/sha1"
	"fmt"
)

func Sha1String(data []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(data))
}
