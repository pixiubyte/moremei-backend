package http

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// MobileRegex 手机号正则表达式 https://github.com/VincentSit/ChinaMobilePhoneNumberRegex/blob/master/README-CN.md
const MobileRegex = "^(?:\\+?86)?1(?:3\\d{3}|5[^4\\D]\\d{2}|8\\d{3}|7(?:[0-35-9]\\d{2}|4(?:0\\d|1[0-2]|9\\d))|9[0-35-9]\\d{2}|6[2567]\\d{2}|4[579]\\d{2})\\d{6}$"

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	if len(mobile) == 0 {
		return false
	}

	match, err := regexp.Match(MobileRegex, []byte(mobile))
	if err != nil {
		return false
	}

	return match
}
