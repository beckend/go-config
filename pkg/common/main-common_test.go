package common_test

import (
	errors "errors"
	testing "testing"

	common "github.com/beckend/go-config/pkg/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPkg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pkg common Suite")
}

var _ = Describe("pkg common", func() {
	Context("FailOnError", func() {
		When("when not nil", func() {
			It("panics", func() {
				Expect(func() { common.FailOnError(errors.New("math: square root of negative number")) }).Should(Panic())
			})
		})
	})
})
