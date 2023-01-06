package config

import "github.com/go-playground/validator/v10"

type DefaultValidator struct {
	Validator *validator.Validate
}

func (dv *DefaultValidator) Validate(i interface{}) error {
	if err := dv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}
