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
	It("can be created by calling New", func() {

		def := New(Describe, true, true, "arg1")
		Expect(def).ShouldNot(BeNil())
	})

	It("can be created by calling New without args", func() {

		def := New(Describe, true, true)
		Expect(def).ShouldNot(BeNil())
		Expect(def.Get().Args).Should(HaveLen(0))
	})

	It("should be able to replace functions args", func() {

		def := New(Describe, true, true, "arg1")
		Expect(def).ShouldNot(BeNil())
		Expect(def.Get().Args).Should(HaveLen(1))
		def.ReplaceArgs("arg1", "arg2")
		Expect(def.Get().Args).Should(HaveLen(2))
	})
})
