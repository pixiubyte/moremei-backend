package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func CheckUriInSlice(uri string, method string, list []string) bool {
	flag := false
	for _, v := range list {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}

		allowedParams := strings.Split(v, ":")
		if len(allowedParams) < 2 {
			continue
		}

		matched, err := regexp.Match(fmt.Sprintf("^%s$", allowedParams[0]), []byte(uri))
		if err != nil {
			continue
		}

		if matched && strings.ToUpper(allowedParams[1]) == strings.ToUpper(method) {
			flag = true
			break
		}
	}

	return flag
}
