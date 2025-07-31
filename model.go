package mox

import "github.com/samber/mo"

type MoUser struct {
	Name mo.Option[string] `validate:"notnil,min=1"`
}
