package singletons_test

import (
	testing "testing"

	singletons "github.com/beckend/go-config/pkg/singletons"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPkgAnagram(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pkg singletons Suite")
}

var _ = Describe("pkg singletons", func() {
	Context("GetSingletons", func() {
		When("called", func() {
			It("will not panic and value is set", func() {
				result := singletons.GetSingletons()
				Expect(result.Validation).ToNot(BeNil())
			})
		})
	})
})
