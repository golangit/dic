package container

import (
	"errors"
	"fmt"
	"github.com/golangit/dic/definition"
	"github.com/golangit/dic/reference"
	"reflect"
	"strings"
)

/*
@todo
- has
- remove
*/

const (
	TagPrefix = "dic"
)

var (
	ErrParamsNotAdapted = func(paramGiven int, paramFun int) error {
		return errors.New(fmt.Sprintf("The number of params is not adapted. Given:%d, Expected:%d", paramGiven, paramFun))
	}

	ErrServiceNotFound       = func(service string) error { return errors.New("Service " + service + " not found") }
	ErrServiceNotPublic      = func(service string) error { return errors.New("Service " + service + " is not public") }
	ErrObjectIsNotAReference = errors.New("Object is not a Resource type")
	ErrKindIsNotValid        = func(kind definition.Kind) error {
		return errors.New("Kind " + string(kind) + " is not a valid kind")
	}
)

type container struct {
	definitions map[string]definition.Definition
	services    map[string][]reflect.Value
}

type Container interface {
	ContainerReader
	ContainerWriter
	ContainerCaller
}

type ContainerReader interface {
	GetDefinition(string) (definition.Definition, error)
}

type ContainerWriter interface {
	Register(string, interface{}, ...interface{}) error
	SetDefinition(string, definition.Definition) error
	Inject(interface{}) error
}

type ContainerCaller interface {
	GetAll(string, ...interface{}) ([]reflect.Value, error)
	Get(string, ...interface{}) interface{}
}

func New() Container {
	return &container{
		definitions: make(map[string]definition.Definition),
		services:    make(map[string][]reflect.Value),
	}
}

func (cnt *container) isFieldTagged(name string) bool {
	return strings.HasPrefix(name, TagPrefix)
}

func (cnt *container) getServiceNameByFieldTagged(name string) string {
	return strings.TrimPrefix(name, TagPrefix)
}

func (cnt *container) Inject(val interface{}) (err error) {

	v := reflect.ValueOf(val)

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil // Should not panic here ?
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		tag := t.Field(i).Tag.Get("dic")
		fmt.Printf("tag [%s], %v", tag, t.Field(i).Tag)
		if f.CanSet() && len(tag) > 0 {
			s := cnt.resolveADependency(tag)
			fmt.Printf("tag %s, %v", tag, s)
			f.Set(s)
		}
	}

	return
}

func (cnt *container) Register(name string, fn interface{}, params ...interface{}) (err error) {

	return cnt.SetDefinition(name, definition.New(fn, true, true, params...))
}

func (cnt *container) SetDefinition(name string, definition definition.Definition) (err error) {

	cnt.definitions[name] = definition
	return
}

func (cnt *container) GetDefinition(name string) (def definition.Definition, err error) {

	if _, ok := cnt.definitions[name]; !ok {
		return nil, ErrServiceNotFound(name)
	}

	return cnt.definitions[name], nil
}

// Get the service, and return only the first result's interface
func (cnt *container) Get(name string, params ...interface{}) interface{} {

	results, _ := cnt.GetAll(name, params...)

	if len(results) < 1 {
		return nil
	}

	return results[0].Interface()
}

// Get a service by its name
// it resolves also the dependencies
// if the service is defined as static it instanciate the object once
func (cnt *container) GetAll(name string, params ...interface{}) ([]reflect.Value, error) {

	def, err := cnt.GetDefinition(name)
	if err != nil {
		return nil, err
	}

	if !def.IsPublic() {
		return nil, ErrServiceNotPublic(name)
	}

	return cnt.doGetAll(name, params...)
}

// call the correct function by its type
func (cnt *container) doGetAll(name string, params ...interface{}) (result []reflect.Value, err error) {

	def, err := cnt.GetDefinition(name)
	if err != nil {
		return
	}
	// cache hit:
	// -if was already stored
	// -and there are no paramters
	// -and the call is static get the already stored value
	if _, ok := cnt.services[name]; ok && len(params) == 0 && def.IsStatic() {
		result = cnt.services[name]
		return
	}

	if len(params) > 0 {
		def.ReplaceArgs(params...)
	}

	dependencies := cnt.resolveDependencies(def.GetArgs())

	switch def.GetKind() {
	case definition.Callable:
		result, err = cnt.callAndInjectACallable(def, dependencies)
	case definition.Fields:
		result, err = cnt.getAndInjectAllFields(def, dependencies)
	case definition.Parameter:
		result = make([]reflect.Value, 1)
		result[0] = def.GetValue()
	default:
		return nil, ErrKindIsNotValid(def.GetKind())
	}

	// store cache:
	// -if no param
	if len(params) == 0 {
		cnt.services[name] = result
	}

	return
}

// inject all the dependency into the object with fields like struct
func (cnt *container) getAndInjectAllFields(def definition.Definition, params []reflect.Value) (result []reflect.Value, err error) {

	v := def.GetValue()

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.CanSet() {
			f.Set(params[i])
		}
	}

	result = make([]reflect.Value, 1)
	result[0] = v
	return
}

func (cnt *container) callAndInjectACallable(def definition.Definition, params []reflect.Value) (result []reflect.Value, err error) {

	// if is a not callable panic here.
	if len(params) != def.GetValue().Type().NumIn() {
		err = ErrParamsNotAdapted(len(params), def.GetValue().Type().NumIn())
		return
	}
	result = def.GetValue().Call(params)
	return
}

func (cnt *container) resolveADependency(name string) reflect.Value {
	ret, _ := cnt.doGetAll(name)
	return ret[0]
}

func (cnt *container) resolveDependencies(params []reflect.Value) []reflect.Value {

	deps := make([]reflect.Value, len(params))

	for k, param := range params {
		if reference, err := IsAReference(param); err == nil {
			deps[k] = cnt.resolveADependency(reference)
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
