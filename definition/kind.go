package definition

import (
	"reflect"
)

type Kind int

const (
	Invalid   Kind = iota
	Callable       // is a callable object, could have params with injection
	Fields         // could have fields with injection like struct, array, maps
	Parameter      // is not callable, and doesn't have fields like string int etc...
)

func GetCorrectKind(k reflect.Kind) Kind {

	switch k {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.Chan,
		reflect.String,
		reflect.Array,
		reflect.Slice,
		reflect.Map:
		return Parameter

	case reflect.Func:
		return Callable

	case reflect.Struct:
		return Fields

	case reflect.Invalid,
		reflect.Uintptr,
		reflect.Interface,
		reflect.Ptr,
		reflect.UnsafePointer:
		return Invalid
	}

	return Invalid
}
