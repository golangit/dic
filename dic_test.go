package dic_test

import (
	"github.com/golangit/dic"
	"reflect"
	"testing"
)

func assertNotEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
func assert(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestNew(t *testing.T) {
	container := dic.New()
	assertNotEqual(t, container, false)
}

func TestRegisterAndGet(t *testing.T) {
	container := dic.New()
	assertion := "Hello"
	key := "key"

	container.Register(key, assertion)
	assertNotEqual(t, container, false)

	value := container.Get(key).String()
	assert(t, assertion, value)
}
