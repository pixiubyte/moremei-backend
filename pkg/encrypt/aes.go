package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

type Aes struct {
	byteKey []byte
}

func NewAes(key string) *Aes {
	l := len(key)
	if l < 32 {
		padLength := 32 - l
		padding := strings.Repeat("b", padLength)
		key = key + padding
	} else {
		key = key[:32]
	}
	b, _ := hex.DecodeString(key)

	return &Aes{
		byteKey: b,
	}

}

func (a *Aes) Encrypt(text []byte) (string, error) {
	block, err := aes.NewCipher(a.byteKey)
	if err != nil {
		return "", err
	}

	// 填充明文数据
	text = pkcs7Pad(text, aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], text)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (a *Aes) Decrypt(data string) ([]byte, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(a.byteKey)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// 去除填充字节
	plaintext, err := pkcs7Unpad(ciphertext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:(length - unpadding)], nil
}
