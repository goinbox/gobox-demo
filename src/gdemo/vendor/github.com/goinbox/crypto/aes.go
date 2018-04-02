package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"strconv"
)

const (
	AES_BLOCK_SIZE = 16
)

type AesCBCCrypter struct {
	blockSize int

	encryptBlockMode cipher.BlockMode
	decryptBlockMode cipher.BlockMode

	padding PaddingInterface
}

func NewAesCBCCrypter(key []byte, iv []byte) (*AesCBCCrypter, error) {
	l := len(key)
	if l != 32 && l != 24 && l != 16 {
		return nil, errors.New("The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.")
	}

	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		return nil, errors.New("The length of iv must be the same as the Block's block size " + strconv.Itoa(blockSize))
	}

	a := &AesCBCCrypter{
		blockSize: blockSize,

		encryptBlockMode: cipher.NewCBCEncrypter(block, iv),
		decryptBlockMode: cipher.NewCBCDecrypter(block, iv),

		padding: &PKCS5Padding{
			BlockSize: blockSize,
		},
	}

	return a, nil
}

func (a *AesCBCCrypter) BlockSize() int {
	return a.blockSize
}

func (a *AesCBCCrypter) SetPadding(padding PaddingInterface) {
	a.padding = padding
}

func (a *AesCBCCrypter) Encrypt(data []byte) []byte {
	data = a.padding.Padding(data)

	crypted := make([]byte, len(data))
	a.encryptBlockMode.CryptBlocks(crypted, data)

	return crypted
}

func (a *AesCBCCrypter) Decrypt(crypted []byte) []byte {
	data := make([]byte, len(crypted))
	a.decryptBlockMode.CryptBlocks(data, crypted)

	return a.padding.UnPadding(data)
}
