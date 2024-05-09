package common

import "github.com/go-playground/validator/v10"

func NewValidator() (v *validator.Validate) {
	v = validator.New(validator.WithRequiredStructEnabled())
	return v
}
