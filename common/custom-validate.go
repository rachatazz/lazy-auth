package common

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func PasswordValidate(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if ok {
		tests := []string{".{8,}", "[a-z]", "[A-Z]", "[0-9]"}
		for _, test := range tests {
			t, _ := regexp.MatchString(test, password)
			if !t {
				return false
			}
		}
	}
	return true
}
