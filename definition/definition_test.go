package definition_test

import (
	. "github.com/golangit/dic/definition"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestDic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Definition Suite")
}

var _ = Describe("Definition/Definition", func() {
	It("can be created by calling New by passing a callable function", func() {

		def := New(Describe, true, true, "arg1")
		Expect(def).ShouldNot(BeNil())
	})

	It("can be created by calling New by passing a struct", func() {
		type TestStruct struct {
			Nerd string
		}
		test := &TestStruct{}
		def := New(test, true, true, "liuggio")
		Expect(def).ShouldNot(BeNil())
	})

	It("can be created by calling New without args", func() {

		def := New(Describe, true, true)
		Expect(def).ShouldNot(BeNil())
		Expect(def.GetArgs()).Should(HaveLen(0))
	})

	It("should be able to replace functions args", func() {

		def := New(Describe, true, true, "arg1")
		Expect(def).ShouldNot(BeNil())
		Expect(def.GetArgs()).Should(HaveLen(1))
		def.ReplaceArgs("arg1", "arg2")
		Expect(def.GetArgs()).Should(HaveLen(2))
	})
})
