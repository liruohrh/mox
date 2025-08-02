package mox

import (
	"github.com/go-playground/validator/v10"
	"github.com/samber/mo"
	"reflect"
)

// RegisterGPVUnwrapOptionTypeFunc why unwrap? because use value for other validate func
func RegisterGPVUnwrapOptionTypeFunc(validate *validator.Validate) {
	validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if v, ok := field.Interface().(mo.Option[string]); ok {
			if v.IsPresent() {
				return v.MustGet()
			} else {
				return nil
			}
		}
		return nil
	}, mo.Option[string]{})
}

// RegisterGPValidatorPresent require option.IsPresent=true
func RegisterGPValidatorPresent(validate *validator.Validate) error {
	return validate.RegisterValidation("present", gpValidatorPresent)
}
func gpValidatorPresent(fl validator.FieldLevel) bool {
	field := fl.Field()
	if !IsOption(field.Type()) {
		return true
	}
	return field.MethodByName("IsPresent").Call(nil)[0].Interface().(bool)
}

// RegisterGPValidatorNotNil notnil: mandatory, allows zero value (except nil)
//
//	required: mandatory, requires non-zero value
//	omitnil: optional, requires non-zero value (except nil)
//	omitempty: optional, allows zero value
func RegisterGPValidatorNotNil(validate *validator.Validate) error {
	return validate.RegisterValidation("notnil", gpValidatorNotNil)
}
func gpValidatorNotNil(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		return true // No validation for non-reference types
	}
}
