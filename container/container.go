package container

import (
	"errors"
	"fmt"
	"github.com/golangit/dic/definition"
	"github.com/golangit/dic/reference"
	"reflect"
)

var (
	ErrParamsNotAdapted = func(paramGiven int, paramFun int) error {
		return errors.New(fmt.Sprintf("The number of params is not adapted. Given:%d, Expected:%d", paramGiven, paramFun))
	}

	ErrServiceNotFound       = func(service string) error { return errors.New("Service " + service + " not found") }
	ErrServiceNotPublic      = func(service string) error { return errors.New("Service " + service + " is not public") }
	ErrObjectIsNotAReference = errors.New("Object is not a Resource type")
)

type container struct {
	definitions map[string]definition.Definition
	services    map[string][]reflect.Value
}

type Container interface {
	ContainerReader
	ContainerWriter
}

type ContainerReader interface {
	GetDefinition(string) definition.Definition
	Call(string, []reflect.Value) ([]reflect.Value, error)
}

type ContainerWriter interface {
	Register(string, interface{}, ...interface{}) error
	GetAll(string, ...interface{}) ([]reflect.Value, error)
	Get(string, ...interface{}) interface{}
	SetDefinition(string, definition.Definition) error
}

func New() Container {
	return &container{
		definitions: make(map[string]definition.Definition),
		services:    make(map[string][]reflect.Value),
	}
}

func (cnt *container) Register(name string, fn interface{}, params ...interface{}) (err error) {

	return cnt.SetDefinition(name, definition.New(fn, true, true, params...))
}

func (cnt *container) SetDefinition(name string, definition definition.Definition) (err error) {

	cnt.definitions[name] = definition
	return
}

func (cnt *container) GetDefinition(name string) definition.Definition {

	return cnt.definitions[name]
}

// Get a service by its name
// it resolves also the dependencies
// if the service is defined as static it instanciate the object once
func (cnt *container) GetAll(name string, params ...interface{}) (result []reflect.Value, err error) {

	if _, ok := cnt.definitions[name]; !ok {
		err = ErrServiceNotFound(name)
		return
	}

	if !cnt.definitions[name].Get().Public {
		err = ErrServiceNotPublic(name)
		return
	}

	// cache hit:
	// -if was already stored
	// -and there are no paramters
	// -and the call is static get the already stored value
	if _, ok := cnt.services[name]; ok && len(params) == 0 && cnt.definitions[name].Get().Static {
		result = cnt.services[name]
		return
	}

	if len(params) > 0 {
		cnt.definitions[name].ReplaceArgs(params...)
	}

	dependencies := cnt.resolveDependencies(cnt.definitions[name].Get().Args)
	result, err = cnt.Call(name, dependencies)

	// store cache:
	// -if no param
	if len(params) == 0 {
		cnt.services[name] = result
	}

	return
}

// Get the service, and return only the first argument's interface
func (cnt *container) Get(name string, params ...interface{}) interface{} {

	results, _ := cnt.GetAll(name, params...)
	return results[0].Interface()
}

func (cnt *container) Call(name string, params []reflect.Value) (result []reflect.Value, err error) {

	if _, ok := cnt.definitions[name]; !ok {
		err = ErrServiceNotFound(name)
		return
	}
	if len(params) != cnt.definitions[name].Get().Value.Type().NumIn() {
		err = ErrParamsNotAdapted(len(params), cnt.definitions[name].Get().Value.Type().NumIn())
		return
	}

	result = cnt.definitions[name].Get().Value.Call(params)
	// todo get the interface value of the first?
	return
}

func (cnt *container) resolveDependencies(params []reflect.Value) []reflect.Value {

	deps := make([]reflect.Value, len(params))

	for k, param := range params {
		if reference, err := IsAReference(param); err == nil {
			ret, _ := cnt.GetAll(reference)
			deps[k] = ret[0]
		} else {
			deps[k] = param
		}
	}
	return deps
}

func getReference(val reference.Reference) string {
	return val.Reference()
}

func IsAReference(value reflect.Value) (ret string, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = ErrObjectIsNotAReference
		}
	}()

	ret = getReference(value.Interface().(reference.Reference))
	return
}
