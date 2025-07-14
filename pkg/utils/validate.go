package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func GetValidate() *validator.Validate {
	return validate
}
