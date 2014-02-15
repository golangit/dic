package reference_test

import (
	. "github.com/golangit/dic/reference"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestReference(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Reference Suite")
}

var _ = Describe("Reference", func() {
	It("can be created by calling New", func() {
		ref := New("ref")
		Expect(ref.Reference()).Should(Equal("ref"))
	})
})
