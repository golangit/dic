package definition

import (
	"reflect"
)

const (
	SERVICE_STATIC  = true
	SERVICE_PUBLIC  = true
	SERVICE_PRIVATE = false
)

type definition struct {
	Value  reflect.Value
	Args   []reflect.Value
	Public bool
	Static bool
}

type Definition interface {
	DefinitionWriter
}

type DefinitionWriter interface {
	ReplaceArgs(params ...interface{}) Definition
	Get() *definition
	SetPublic(bool) Definition
	SetStatic(bool) Definition
}

func New(fn interface{}, pub bool, stat bool, params ...interface{}) Definition {
	// todo change error
	v := reflect.ValueOf(fn)
	//v.Type().NumIn() // panic if is not callable

	var in = make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	return &definition{
		Value:  v,
		Args:   in,
		Public: pub,
		Static: stat,
	}
}

// add arguments to a defined Definition
func (d *definition) ReplaceArgs(params ...interface{}) Definition {

	var in = make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	d.Args = in

	return d
}

func (d *definition) SetStatic(isStatic bool) Definition {

	d.Static = isStatic

	return d
}

func (d *definition) SetPublic(isPublic bool) Definition {

	d.Public = isPublic

	return d
}

func (d *definition) Get() *definition {

	return d
}
