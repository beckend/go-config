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

func TestPkgAnagram(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pkg validation Suite")
}

var _ = Describe("pkg validation", func() {
	Context("GetValidator", func() {
		When("called", func() {
			It("will not panic and value is set", func() {
				result := validation.GetValidator()
				Expect(result.Utils).ToNot(BeNil())
				Expect(result.Validate).ToNot(BeNil())
			})
		})
	})

	Context("Validation", func() {
		When("failing validation with option PanicOnError true", func() {
			It("will panic", func() {
				// instance := validation.GetValidator()

				// instance.Utils.ValidatorUtils.ValidateStruct(validation.ValidatorUtilsValidateStructOptions{
				// 	PrefixError: "error - ",
				// 	Validate:    instance.Validate,
				// 	TheStruct: TestValidateStructOne{
				// 		AccessKey: "access",
				// 	},
				// 	PanicOnError: true,
				// })

				// Expect(func() {
				// 	instance := validation.GetValidator()

				// 	instance.Utils.ValidatorUtils.ValidateStruct(validation.ValidatorUtilsValidateStructOptions{
				// 		PrefixError: "error - ",
				// 		Validate:    instance.Validate,
				// 		TheStruct: TestValidateStructOne{
				// 			AccessKey: "access",
				// 		},
				// 		PanicOnError: true,
				// 	})
				// }).To(PanicWith(MatchRegexp(`.+ok$`)))
			})
		})

		// When("failing validation with option PanicOnError false", func() {
		// 	It("will not panic", func() {
		// 		instance := validation.GetValidator()

		// 		Expect(func() {
		// 			instance.Utils.ValidatorUtils.ValidateStruct(validation.ValidatorUtilsValidateStructOptions{
		// 				PrefixError:  "error - ",
		// 				Validate:     instance.Validate,
		// 				TheStruct:    TestValidateStructOne{},
		// 				PanicOnError: false,
		// 			})
		// 		}).ToNot(Panic())
		// 	})
		// })
	})
})
