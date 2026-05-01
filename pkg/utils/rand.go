package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func RandCode(len int) string {
	code := ""

	rand.New(rand.NewSource(time.Now().Unix())).Seed(time.Now().Unix())
	for i := 0; i < len; i++ {
		if i == 0 {
			code = fmt.Sprintf("%s%d", code, rand.Intn(9)+1)
		} else {
			code = fmt.Sprintf("%s%d", code, rand.Intn(10))
		}
	}

	return code
}
