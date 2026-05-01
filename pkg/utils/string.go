package utils

import "regexp"

func IsNumeric(s string) bool {
	regex := regexp.MustCompile(`^\d+$`)
	return regex.MatchString(s)
}
