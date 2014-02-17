package definition

import (
	"reflect"
)

type definition struct {
	value  reflect.Value
	args   []reflect.Value
	public bool
	static bool
	kind   Kind
}

type Definition interface {
	DefinitionWriter
	DefinitionReader
}

type DefinitionReader interface {
	GetValue() reflect.Value
	GetArgs() []reflect.Value
	IsPublic() bool
	IsStatic() bool
	GetKind() Kind
}

type DefinitionWriter interface {
	ReplaceArgs(params ...interface{}) Definition
	Get() *definition
	SetPublic(bool) Definition
	SetStatic(bool) Definition
}

func New(fn interface{}, pub bool, stat bool, params ...interface{}) Definition {

	v := reflect.ValueOf(fn)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var in = make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	return &definition{
		value:  v,
		args:   in,
		public: pub,
		static: stat,
		kind:   GetCorrectKind(v.Kind()),
	}
}

// add arguments to a defined Definition
func (d *definition) ReplaceArgs(params ...interface{}) Definition {

	var in = make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	d.args = in
	return d
}

func (d *definition) SetStatic(isStatic bool) Definition {
	d.static = isStatic
	return d
}

func (d *definition) SetPublic(isPublic bool) Definition {
	d.public = isPublic
	return d
}

func (d *definition) Get() *definition {
	return d
}

func (d *definition) GetValue() reflect.Value {
	return d.value
}

func (d *definition) GetArgs() []reflect.Value {
	return d.args
}

func (d *definition) IsStatic() bool {
	return d.static
}

func (d *definition) IsPublic() bool {
	return d.public
}

func (d *definition) GetKind() Kind {
	return d.kind
}
