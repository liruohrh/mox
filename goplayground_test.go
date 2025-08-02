package mox

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/samber/mo"
	"github.com/stretchr/testify/require"
)

func TestGoPlaygroundPresent(t *testing.T) {
	validate := validator.New()
	require.NoError(t, RegisterGPValidatorPresent(validate))
	datas := []struct {
		value    ValidatePresentDto
		contains string
	}{
		{
			value: ValidatePresentDto{
				V: mo.Some(""),
			},
			contains: "",
		},
		{
			value: ValidatePresentDto{
				V: mo.None[string](),
			},
			contains: "failed on the 'present' tag",
		},
	}

	for _, data := range datas {
		err := validate.Struct(&data.value)
		if data.contains == "" {
			require.NoError(t, err)
		} else {
			require.ErrorContains(t, err, data.contains)
		}
	}
}

func TestGoPlaygroundNotnil(t *testing.T) {
	validate := validator.New()
	require.NoError(t, RegisterGPValidatorNotNil(validate))
	datas := []struct {
		value    ValidateNotNilOptionDto
		contains string
	}{
		{
			value: ValidateNotNilOptionDto{
				V: mo.Some(""),
			},
			contains: "",
		},
		{
			value:    ValidateNotNilOptionDto{},
			contains: "failed on the 'notnil' tag",
		},
		{
			value: ValidateNotNilOptionDto{
				V: mo.None[string](),
			},
			contains: "failed on the 'notnil' tag",
		},
		{
			value: ValidateNotNilOptionDto{
				V: mo.None[string](),
			},
			contains: "failed on the 'notnil' tag",
		},
	}

	for _, data := range datas {
		err := validate.Struct(&data.value)
		if data.contains == "" {
			require.NoError(t, err)
		} else {
			require.ErrorContains(t, err, data.contains)
		}
	}

	str := ""
	datas2 := []struct {
		value    ValidateNotNilPointerDto
		contains string
	}{
		{
			value: ValidateNotNilPointerDto{
				V: &str,
			},
			contains: "",
		},
		{
			value:    ValidateNotNilPointerDto{},
			contains: "failed on the 'notnil' tag",
		},
	}

	for _, data := range datas2 {
		err := validate.Struct(&data.value)
		if data.contains == "" {
			require.NoError(t, err)
		} else {
			require.ErrorContains(t, err, data.contains)
		}
	}
}

func TestGoPlaygroundOptionValidate(t *testing.T) {
	validate := validator.New()
	RegisterGPVUnwrapOptionTypeFunc(validate)
	datas := []struct {
		value    ValidateOptionDto
		contains string
	}{
		{
			value:    ValidateOptionDto{},
			contains: "failed on the 'min' tag",
		},
		{
			value: ValidateOptionDto{
				V: mo.None[string](),
			},
			contains: "failed on the 'min' tag",
		},
		{
			value: ValidateOptionDto{
				V: mo.Some("1234"),
			},
			contains: "failed on the 'min' tag",
		},
		{
			value: ValidateOptionDto{
				V: mo.Some("12345"),
			},
			contains: "",
		},
	}

	for _, data := range datas {
		err := validate.Struct(&data.value)
		if data.contains == "" {
			require.NoError(t, err)
		} else {
			require.ErrorContains(t, err, data.contains)
		}
	}
}

func TestGoPlaygroundValidator(t *testing.T) {
	validatorI := validator.New()

	emptyStr := ""

	// No ignore setting, so min validation fails
	require.ErrorContains(t, validatorI.Struct(&struct {
		V *string `validate:"min=5"`
	}{}), "failed on the 'min' tag")
	// Set to ignore nil
	require.NoError(t, validatorI.Struct(&struct {
		V *string `validate:"omitnil,min=5"`
	}{}))

	// Set required to require non-zero value
	require.ErrorContains(t, validatorI.Struct(&struct {
		V *string `validate:"required,min=5"`
	}{}), "failed on the 'required' tag")
	// required is based on type, for pointers it only requires non-nil
	require.ErrorContains(t, validatorI.Struct(&struct {
		V *string `validate:"required,min=5"`
	}{
		V: &emptyStr,
	}), "failed on the 'min' tag")

	// Set omitempty to ignore zero values, opposite of required
	require.NoError(t, validatorI.Struct(&struct {
		V *string `validate:"omitempty,min=5"`
	}{}))
	require.ErrorContains(t, validatorI.Struct(&struct {
		V *string `validate:"omitempty,min=5"`
	}{
		V: &emptyStr,
	}), "failed on the 'min' tag")
}
