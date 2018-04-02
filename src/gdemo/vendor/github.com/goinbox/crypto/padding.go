package crypto

import (
	"bytes"
)

type PaddingInterface interface {
	Padding(data []byte) []byte
	UnPadding(data []byte) []byte
}

type PKCS5Padding struct {
	BlockSize int
}

func (p *PKCS5Padding) Padding(data []byte) []byte {
	padding := p.BlockSize - len(data)%p.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(data, padtext...)
}

func (p *PKCS5Padding) UnPadding(data []byte) []byte {
	l := len(data)
	unpadding := int(data[l-1])

	return data[:(l - unpadding)]
}
