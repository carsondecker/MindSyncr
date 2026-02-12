package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func RegisterCustomValidations(v *validator.Validate) {
	v.RegisterValidation("not_nil_uuid", notNilUUID)
}

func notNilUUID(fl validator.FieldLevel) bool {
	u, ok := fl.Field().Interface().(uuid.UUID)
	if !ok {
		return false
	}
	return u != uuid.Nil
}
