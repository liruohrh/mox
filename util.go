package mox

import (
	"reflect"
	"strings"
)

func IsOption(ot reflect.Type) bool {
	pkg := ot.PkgPath()
	typeName := ot.Name()
	if pkg == "github.com/samber/mo" && strings.HasPrefix(typeName, "Option[") {
		return true
	}
	return false
}
