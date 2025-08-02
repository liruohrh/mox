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
type GinDto struct {
	// Basic types
	Str      string  `form:"str"`
	Bool     bool    `form:"bool"`
	Intv     int     `form:"intv"`
	Int8v    int8    `form:"int8v"`
	Int16v   int16   `form:"int16v"`
	Int32v   int32   `form:"int32v"`
	Int64v   int64   `form:"int64v"`
	UIntv    uint    `form:"uintv"`
	UInt8v   uint8   `form:"uint8v"`
	UInt16v  uint16  `form:"uint16v"`
	UInt32v  uint32  `form:"uint32v"`
	UInt64v  uint64  `form:"uint64v"`
	Float32V float32 `form:"float32V"`
	Float64V float64 `form:"float64V"`

	// Pointer types
	PtrStr     *string  `form:"ptrStr"`
	PtrBool    *bool    `form:"ptrBool"`
	PtrInt     *int     `form:"ptrInt"`
	PtrInt8    *int8    `form:"ptrInt8"`
	PtrInt16   *int16   `form:"ptrInt16"`
	PtrInt32   *int32   `form:"ptrInt32"`
	PtrInt64   *int64   `form:"ptrInt64"`
	PtrUInt    *uint    `form:"ptrUInt"`
	PtrUInt8   *uint8   `form:"ptrUInt8"`
	PtrUInt16  *uint16  `form:"ptrUInt16"`
	PtrUInt32  *uint32  `form:"ptrUInt32"`
	PtrUInt64  *uint64  `form:"ptrUInt64"`
	PtrFloat32 *float32 `form:"ptrFloat32"`
	PtrFloat64 *float64 `form:"ptrFloat64"`

	// Slice types
	SliceStr     []string  `form:"sliceStr"`
	SliceBool    []bool    `form:"sliceBool"`
	SliceInt     []int     `form:"sliceInt"`
	SliceInt8    []int8    `form:"sliceInt8"`
	SliceInt16   []int16   `form:"sliceInt16"`
	SliceInt32   []int32   `form:"sliceInt32"`
	SliceInt64   []int64   `form:"sliceInt64"`
	SliceUInt    []uint    `form:"sliceUInt"`
	SliceUInt8   []uint8   `form:"sliceUInt8"`
	SliceUInt16  []uint16  `form:"sliceUInt16"`
	SliceUInt32  []uint32  `form:"sliceUInt32"`
	SliceUInt64  []uint64  `form:"sliceUInt64"`
	SliceFloat32 []float32 `form:"sliceFloat32"`
	SliceFloat64 []float64 `form:"sliceFloat64"`

	// Array types
	ArrayStr     [2]string  `form:"arrayStr"`
	ArrayBool    [2]bool    `form:"arrayBool"`
	ArrayInt     [2]int     `form:"arrayInt"`
	ArrayInt8    [2]int8    `form:"arrayInt8"`
	ArrayInt16   [2]int16   `form:"arrayInt16"`
	ArrayInt32   [2]int32   `form:"arrayInt32"`
	ArrayInt64   [2]int64   `form:"arrayInt64"`
	ArrayUInt    [2]uint    `form:"arrayUInt"`
	ArrayUInt8   [2]uint8   `form:"arrayUInt8"`
	ArrayUInt16  [2]uint16  `form:"arrayUInt16"`
	ArrayUInt32  [2]uint32  `form:"arrayUInt32"`
	ArrayUInt64  [2]uint64  `form:"arrayUInt64"`
	ArrayFloat32 [2]float32 `form:"arrayFloat32"`
	ArrayFloat64 [2]float64 `form:"arrayFloat64"`

	// Option types - Basic
	StrOption     mo.Option[string]  `form:"strOption"`
	BoolOption    mo.Option[bool]    `form:"boolOption"`
	IntOption     mo.Option[int]     `form:"intOption"`
	Int8Option    mo.Option[int8]    `form:"int8Option"`
	Int16Option   mo.Option[int16]   `form:"int16Option"`
	Int32Option   mo.Option[int32]   `form:"int32Option"`
	Int64Option   mo.Option[int64]   `form:"int64Option"`
	UIntOption    mo.Option[uint]    `form:"uintOption"`
	UInt8Option   mo.Option[uint8]   `form:"uint8Option"`
	UInt16Option  mo.Option[uint16]  `form:"uint16Option"`
	UInt32Option  mo.Option[uint32]  `form:"uint32Option"`
	UInt64Option  mo.Option[uint64]  `form:"uint64Option"`
	Float32Option mo.Option[float32] `form:"float32Option"`
	Float64Option mo.Option[float64] `form:"float64Option"`

	// Option types - Slice
	SliceStrOption     mo.Option[[]string]  `form:"sliceStrOption"`
	SliceBoolOption    mo.Option[[]bool]    `form:"sliceBoolOption"`
	SliceIntOption     mo.Option[[]int]     `form:"sliceIntOption"`
	SliceInt8Option    mo.Option[[]int8]    `form:"sliceInt8Option"`
	SliceInt16Option   mo.Option[[]int16]   `form:"sliceInt16Option"`
	SliceInt32Option   mo.Option[[]int32]   `form:"sliceInt32Option"`
	SliceInt64Option   mo.Option[[]int64]   `form:"sliceInt64Option"`
	SliceUIntOption    mo.Option[[]uint]    `form:"sliceUIntOption"`
	SliceUInt8Option   mo.Option[[]uint8]   `form:"sliceUInt8Option"`
	SliceUInt16Option  mo.Option[[]uint16]  `form:"sliceUInt16Option"`
	SliceUInt32Option  mo.Option[[]uint32]  `form:"sliceUInt32Option"`
	SliceUInt64Option  mo.Option[[]uint64]  `form:"sliceUInt64Option"`
	SliceFloat32Option mo.Option[[]float32] `form:"sliceFloat32Option"`
	SliceFloat64Option mo.Option[[]float64] `form:"sliceFloat64Option"`
}
