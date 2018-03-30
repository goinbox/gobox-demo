package encoding

import (
	"encoding/base64"
)

func Base64Encode(data []byte) []byte {
	coded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(coded, data)

	return coded
}

func Base64Decode(coded []byte) []byte {
	data := make([]byte, base64.StdEncoding.DecodedLen(len(coded)))
	n, _ := base64.StdEncoding.Decode(data, coded)

	return data[:n]
}

func SlashDecode(coded []byte) []byte {
	cnt := len(coded)
	tbs := make([]byte, cnt)

	j := 0
	for i := 0; i < cnt; i++ {
		if coded[i] == '\\' {
			next := i + 1
			if coded[next] == 'n' {
				tbs[j] = '\n'
				i++
			} else if coded[next] == '\\' {
				tbs[j] = '\\'
				i++
			} else {
				tbs[j] = '\\'
			}
		} else {
			tbs[j] = coded[i]
		}
		j++
	}

	return tbs[0:j]
}
