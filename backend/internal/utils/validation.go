package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	hasOnlyValid = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]+$`)
	hasSymbol    = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
	hasUppercase = regexp.MustCompile(`[A-Z]`)
	hasLowercase = regexp.MustCompile(`[a-z]`)
	hasDigit     = regexp.MustCompile(`[0-9]`)
)

func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("password", validatePassword)
}

func validatePassword(fl validator.FieldLevel) bool {
	str := fl.Field().String()

	if len(str) < 10 || len(str) > 400 {
		return false
	}

	if !hasLowercase.MatchString(str) {
		return false
	}

	if !hasUppercase.MatchString(str) {
		return false
	}

	if !hasOnlyValid.MatchString(str) {
		return false
	}

	if !hasSymbol.MatchString(str) {
		return false
	}

	return true
}
