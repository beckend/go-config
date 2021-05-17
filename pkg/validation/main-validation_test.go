package validation_test

import (
	testing "testing"

	validation "github.com/beckend/go-config/pkg/validation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestValidateStructOne struct {
	AccessKey string `validate:"required"`
}

func TestPkg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pkg validation Suite")
}

var _ = Describe("pkg validation", func() {
	Context("New", func() {
		When("called", func() {
			It("will not panic and value is set", func() {
				result := validation.New()
				Expect(result.Validate).ToNot(BeNil())
			})
		})
	})

	Context("Validation", func() {
		When("valid date", func() {
			It("will not panic", func() {
				Expect(func() {
					validation.New().ValidateStruct(validation.ValidatorUtilsValidateStructOptions{
						PrefixError: "error - ",
						TheStruct: TestValidateStructOne{
							AccessKey: "dsadsad",
						},
						PanicOnError: true,
					})
				}).To(Not(Panic()))
			})
		})

		When("failing validation with option PanicOnError true", func() {
			It("will panic", func() {
				Expect(func() {
					validation.New().ValidateStruct(validation.ValidatorUtilsValidateStructOptions{
						PrefixError: "error - ",
						TheStruct: TestValidateStructOne{
							AccessKey: "",
						},
						PanicOnError: true,
					})
				}).To(PanicWith(MatchRegexp(`.+failed.+`)))
			})
		})

		When("failing validation with option PanicOnError false", func() {
			It("will not panic", func() {
				Expect(func() {
					validation.New().ValidateStruct(validation.ValidatorUtilsValidateStructOptions{
						PrefixError:  "error - ",
						TheStruct:    TestValidateStructOne{},
						PanicOnError: false,
					})
				}).ToNot(Panic())
			})
		})
	})
})
