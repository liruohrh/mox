package mox

import "github.com/samber/mo"

type ValidatePresentDto struct {
	V mo.Option[string] `validate:"present"`
}
type ValidateNotNilOptionDto struct {
	V mo.Option[string] `validate:"notnil"`
}
type ValidateNotNilPointerDto struct {
	V *string `validate:"notnil"`
}
type ValidateOptionDto struct {
	V mo.Option[string] `validate:"min=5"`
}
