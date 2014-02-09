package dic

import (
	"reflect"
)

type Dic interface {
	DicWriter
	DicReader
}

type DicWriter interface {
	Register(string, interface{}) Dic
}

type DicReader interface {
	Get(key string) reflect.Value
}

type dic struct {
	container map[string]reflect.Value
}

// factory of Dic
func New() Dic {
	return &dic{
		container: make(map[string]reflect.Value),
	}
}

// Register a new value
func (d *dic) Register(key string, val interface{}) Dic {
	d.container[key] = reflect.ValueOf(val)
	return d
}

func (d *dic) Get(key string) reflect.Value {
	return d.container[key]
}

// Invoke attempts to call the interface{} provided as a function,
// providing dependencies for function arguments based on Type.
// Returns a slice of reflect.Value representing the returned values of the function.
// Returns an error if the injection fails.
// It panics if f is not a function
/*func (d *dic) Get(key string, params ...interface{}) (result []reflect.Value, err error) {

	value := reflect.ValueOf(i.container[key])

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	return f.Call(), nil
}*/
