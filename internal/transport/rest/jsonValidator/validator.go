package jsonValidator

import (
	"fmt"
	"github.com/go-playground/validator"
	"go.uber.org/multierr"
	"reflect"
)

type JsonValidator struct {
	validator   *validator.Validate
	validations map[string]ValidationFunc
}

func New() *JsonValidator {
	return &JsonValidator{
		validator:   validator.New(),
		validations: make(map[string]ValidationFunc),
	}
}

func (v *JsonValidator) AddValidation(fns ...ValidationFunc) {
	for _, fn := range fns {
		v.validator.RegisterValidation(fn.Name, fn.Fn)
		v.validations[fn.Name] = fn
	}
}

func (v *JsonValidator) Validate(s any) (result error) {
	errs := v.validator.Struct(s)

	for _, err := range errs.(validator.ValidationErrors) {
		field, _ := reflect.TypeOf(s).FieldByName(err.Field())
		jsonTag := field.Tag.Get("json")

		if jsonTag == "" {
			//todo
		}

		err := fmt.Errorf("%s %w", jsonTag, v.validations[err.Tag()].Err)
		result = multierr.Append(result, err)
	}

	return result
}

type ValidationFunc struct {
	Name string
	Fn   func(fl validator.FieldLevel) bool
	Err  error
}
