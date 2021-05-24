package environment_test

import (
	testing "testing"

	environment "github.com/beckend/go-config/pkg/environment"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPkg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pkg environment Suite")
}

var _ = Describe("pkg environment", func() {
	Context("GetEnv", func() {
		When("env exists", func() {
			It("returns correct value", func() {
				Expect(len(environment.GetEnv("SHELL", "")) > 0).To(Equal(true))
			})
		})

		When("env does not exist", func() {
			It("returns fallback", func() {
				Expect(environment.GetEnv("FDSCCVB##csdas#@!CS", "")).To(Equal(""))
			})
		})
	})
})
