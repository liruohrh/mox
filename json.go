// Package mox
//
//	for https://github.com/samber/mo/pull/65 trust set null as set a value
package mox

import (
	"reflect"
	"strings"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

type OptionExtension struct {
	jsoniter.DummyExtension
}

func (ext *OptionExtension) CreateEncoder(typ reflect2.Type) jsoniter.ValEncoder {
	typeName := typ.String()
	if strings.HasPrefix(typeName, "mo.Option[") {
		return &OptionEncoder{typ: typ}
	}
	return nil
}

type OptionEncoder struct {
	typ reflect2.Type
}

func (encoder *OptionEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	obj := encoder.typ.UnsafeIndirect(ptr)
	reflectValue := reflect.ValueOf(obj)
	// when IsEmpty=true, will not call Encode, so directly call get here.
	orEmptyMethod := reflectValue.MethodByName("OrEmpty")
	value := orEmptyMethod.Call([]reflect.Value{})[0].Interface()
	stream.WriteVal(value)
}

func (encoder *OptionEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	obj := encoder.typ.UnsafeIndirect(ptr)
	reflectValue := reflect.ValueOf(obj)
	isPresentMethod := reflectValue.MethodByName("IsPresent")
	isPresent := isPresentMethod.Call(nil)[0].Bool()
	return !isPresent
}
