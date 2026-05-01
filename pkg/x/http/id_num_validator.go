package http

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const IdNumRegex18 = "^[1-9]\\d{5}(18|19|([23]\\d))\\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$"
const IdNumRegex15 = "^[1-9]\\d{5}\\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\\d{2}$"

func ValidateIdNum(fl validator.FieldLevel) bool {
	idNumb := fl.Field().String()
	if len(idNumb) == 0 {
		return false
	}

	var regex string
	if len(idNumb) == 18 {
		regex = IdNumRegex18
	} else if len(idNumb) == 15 {
		regex = IdNumRegex15
	} else {
		return false
	}

	match, err := regexp.Match(regex, []byte(idNumb))
	if err != nil {
		return false
	}

	return match
}
