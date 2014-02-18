package container_test

import (
	"fmt"
	. "github.com/golangit/dic/container"
	"github.com/golangit/dic/definition"
	"github.com/golangit/dic/reference"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
	"testing"
	"time"
)

func TestDic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Container Suite")
}

var _ = Describe("Container", func() {
	It("should create a new instance", func() {
		cnt := New()
		Expect(cnt).ShouldNot(BeNil())
	})

	var (
		goodProviders = map[string]interface{}{
			"hello":  func() { print("hello") },
			"foobar": func(a, b, c int) int { return a + b + c },
		}
	)

	Context("When I register a function", func() {
		It("Should not return nil", func() {
			cnt := New()
			for k, v := range goodProviders {
				Expect(cnt.Register(k, v)).Should(BeNil())
			}
		})
		PIt("should understand when there's a Circular reference", func() {
			cnt := New()
			cnt.Register("orm", func(dbName string) string { return "db:" + dbName })
			Expect(cnt.Get("orm", reference.New("orm"))).Should(BeNil())
		})
	})

	Context("When I try to inject all the dependency", func() {
		It("inject all the dependency for a struct", func() {
			type Orm struct {
				DbTable string `dic:"table"`
			}
			cnt := New()
			cnt.Register("table", "table-name")
			orm := new(Orm)
			cnt.Inject(&orm)

			Expect(orm.DbTable).Should(Equal("table-name"))
		})
	})

	Context("When I try to get a definition", func() {
		It("should return nil if the service doesn't exist", func() {
			cnt := New()
			def, _ := cnt.GetDefinition("hello")
			Expect(def).Should(BeNil())
		})
		It("should return a definition", func() {
			def := definition.New(fmt.Printf, true, true)
			cnt := New()
			cnt.Register("hello", func() { fmt.Printf("hello") })
			Expect(cnt.GetDefinition("hello")).Should(BeAssignableToTypeOf(def))
		})
	})

	Context("When I try to get a service", func() {
		It("should execute it and return a list of result", func() {
			cnt := New()
			cnt.Register("root", func() int { return 2 * 2 })

			def, err := cnt.GetAll("root")

			var i int = def[0].Interface().(int)
			Expect(i).Should(Equal(4))
			Expect(err).Should(BeNil())
		})
		It("should execute it and return the interface of the element", func() {
			cnt := New()
			cnt.Register("root", func() int { return 2 * 2 })

			Expect(cnt.Get("root").(int)).Should(Equal(4))
		})
		It("should not execute if is not public", func() {
			cnt := New()
			def := definition.New(func() int { return 2 * 2 }, false, false)
			cnt.SetDefinition("root", def)
			result, err := cnt.GetAll("root")

			Expect(result).Should(BeNil())
			Expect(err.Error()).Should(ContainSubstring("not public"))
		})

		It("should return err if the service doesn't exist", func() {
			cnt := New()
			def, err := cnt.GetAll("hello")
			Expect(def).Should(BeNil())
			Expect(err.Error()).Should(ContainSubstring("Service hello not found"))
		})
	})

	Context("When I try to get a static service", func() {
		It("should not be executed twice", func() {
			cnt := New()
			cnt.Register("root", time.Now)
			before, err := cnt.GetAll("root")
			Expect(err).Should(BeNil())
			now, err := cnt.GetAll("root")
			Expect(err).Should(BeNil())
			Expect(now).Should(Equal(before))
		})
	})

	Context("When I try to inject a dependency into a callable service", func() {
		It("should resolve a dependency to another service", func() {
			cnt := New()
			cnt.Register("service.4", func() int { return 4 })
			cnt.Register("multiplicator", func(a int) int { return a * 2 })

			val, err := cnt.GetAll("multiplicator", reference.New("service.4"))
			Expect(err).Should(BeNil())
			var i int = val[0].Interface().(int)
			Expect(i).Should(Equal(8))
		})
		It("should resolve a dependency to string", func() {
			cnt := New()
			cnt.Register("db.name", "golang")
			cnt.Register("orm", func(dbName string) string { return "db:" + dbName })

			Expect(cnt.Get("orm", reference.New("db.name"))).Should(Equal("db:golang"))
		})
		It("should resolve a mixed type dependencies", func() {
			cnt := New()
			cnt.Register("service.2", func() int { return 2 })
			cnt.Register("service.sum", func(a int, b int, str string) string { return fmt.Sprintf(str, a, b) })

			val, err := cnt.GetAll("service.sum", reference.New("service.2"), 3, "a:%d,b:%d")
			Expect(err).Should(BeNil())
			var i string = val[0].Interface().(string)
			Expect(i).Should(Equal("a:2,b:3"))
		})
	})

	Context("When I try to understand if a parameter is a Reference", func() {
		It("should return the reference Name if the param is a reference", func() {
			ref := reflect.ValueOf(reference.New("ref"))
			ret, err := IsAReference(ref)
			Expect(err).Should(BeNil())
			Expect(ret).Should(Equal("ref"))
		})
		It("should return err if the param is not a ref", func() {
			ref := reflect.ValueOf(true)
			ret, err := IsAReference(ref)
			Expect(ret).Should(BeEmpty())
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("When I try to get a not static service", func() {
		It("should be executed each time", func() {
			cnt := New()
			def := definition.New(time.Now, true, false)
			cnt.SetDefinition("root", def)
			before, err := cnt.GetAll("root")
			Expect(err).Should(BeNil())
			now, err := cnt.GetAll("root")
			Expect(err).Should(BeNil())
			Expect(now).ShouldNot(Equal(before))
		})
	})

	Context("When I try to Register an objects with fields", func() {
		Context("Using a struct", func() {
			It("should register correctly", func() {
				type TestStruct struct {
					Nerd string
				}
				test := &TestStruct{}
				cnt := New()
				err := cnt.Register("root", test)
				Expect(err).Should(BeNil())
			})

			It("should register it and inject its fields", func() {
				type TestStruct struct {
					Nerd string
				}
				test := &TestStruct{}
				cnt := New()
				err := cnt.Register("struct", test, "liuggio")
				Expect(err).Should(BeNil())

				res, err := cnt.GetAll("struct")
				Expect(res).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
				Expect(res[0].Interface().(TestStruct).Nerd).Should(Equal("liuggio"))
			})
		})
	})

	Context("When I register a parameters", func() {
		It("should register it", func() {
			db := "golang"
			cnt := New()
			err := cnt.Register("database.name", db)
			Expect(err).Should(BeNil())
		})
		It("should register, and get it", func() {
			db := "golang"
			cnt := New()
			err := cnt.Register("database.name", db)
			Expect(err).Should(BeNil())
			res, err := cnt.GetAll("database.name")
			Expect(res).ShouldNot(BeNil())
			Expect(err).Should(BeNil())
			Expect(res[0].Interface().(string)).Should(Equal("golang"))
		})
	})

	PContext("When I register a struct", func() {
		PIt("should inject by tagging", func() {
		})
	})
})
