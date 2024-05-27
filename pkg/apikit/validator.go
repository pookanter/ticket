package apikit

import "github.com/go-playground/validator/v10"

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
