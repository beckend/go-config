package reflection_test

import (
	testing "testing"

	reflection "github.com/beckend/go-config/pkg/reflection"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPkg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pkg reflection Suite")
}

var _ = Describe("pkg reflection", func() {
	Context("GetType", func() {
		When("named struct", func() {
			It("pointer outputs correct name", func() {
				type OTest struct{}
				Expect(reflection.GetType(&OTest{})).To(Equal("*OTest"))
			})

			It("non pointer outputs correct name", func() {
				type OTest struct{}
				Expect(reflection.GetType(OTest{})).To(Equal("OTest"))
			})
		})

		When("unnamed non struct", func() {
			It("pointer outputs correct name", func() {
				var input map[string]interface{}
				Expect(reflection.GetType(&input)).To(Equal("*"))
			})

			It("unnamed non struct non pointer outputs correct name", func() {
				var input map[string]interface{}
				Expect(reflection.GetType(input)).To(Equal(""))
			})
		})
	})

	Context("HasElement", func() {
		It("works for slices", func() {
			collection := []string{"hello", "my", "friend"}
			Expect(reflection.HasElement(collection, "friend")).To(Equal(true))
			Expect(reflection.HasElement(collection, "apple")).To(Equal(false))
		})
	})
})
