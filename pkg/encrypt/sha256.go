package encrypt

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func Sha256(str string) string {
	data := base64.StdEncoding.EncodeToString([]byte(str))

	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func IsSha256(str, hashStr string) bool {
	return Sha256(str) == hashStr
}
