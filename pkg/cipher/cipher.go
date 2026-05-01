package cipher

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/zeromicro/go-zero/core/hash"
)

type (
	Cipher struct {
		category string
		iv       []byte
		block    cipher.Block
	}

	EncryptResult struct {
		ciphertext string
		hash       string
		mask       string
	}
)

func NewCipher(category string, key string, iv string) (c *Cipher, err error) {
	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	return &Cipher{
		category: category,
		iv:       []byte(iv),
		block:    block,
	}, nil
}

func (c *Cipher) Encrypt(plaintext string) *EncryptResult {
	paddingText := pkcs7Padding([]byte(plaintext), c.block.BlockSize())

	blockMode := cipher.NewCBCEncrypter(c.block, c.iv)

	ciphertext := make([]byte, len(paddingText))
	blockMode.CryptBlocks(ciphertext, paddingText)

	ct := base64.StdEncoding.EncodeToString(ciphertext)
	return &EncryptResult{
		ciphertext: ct,
		hash:       hash.Md5Hex([]byte(ct)),
		mask:       Masker.Mask(plaintext, Category(c.category)),
	}
}

func (c *Cipher) Decrypt(ciphertext string) (string, error) {
	src, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(c.block, c.iv)

	paddingText := make([]byte, len(src))
	blockMode.CryptBlocks(paddingText, src)

	plaintext := pkcs7UnPadding(paddingText)

	return string(plaintext), nil
}

func (er *EncryptResult) ToString() string {
	return er.ciphertext
}

func (er *EncryptResult) Mask() string {
	return er.mask
}

func (er *EncryptResult) Hash() string {
	return er.hash
}

func pkcs7Padding(unPaddingText []byte, blockSize int) []byte {
	padding := blockSize - len(unPaddingText)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(unPaddingText, paddingText...)
}

func pkcs7UnPadding(paddingText []byte) []byte {
	length := len(paddingText)
	unPadding := int(paddingText[length-1])

	return paddingText[:(length - unPadding)]
}
