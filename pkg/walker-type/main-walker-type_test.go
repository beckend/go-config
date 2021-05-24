package walkertype_test

import (
	reflect "reflect"
	testing "testing"

	walkertype "github.com/beckend/go-config/pkg/walker-type"
	jsoniter "github.com/json-iterator/go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func TestPkg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pkg walkertype Suite")
}

var _ = Describe("pkg walkertype", func() {
	Context("Walk", func() {
		When("walking maps", func() {
			It("recursive loops", func() {
				input := walkertype.TypeMap{
					"key1": walkertype.TypeMap{
						"key2": walkertype.TypeMap{
							"key3": walkertype.TypeArray{"arrayval1", "arrayval2", "arrayval3"},
						},
						"key4": "string1",
						"key5": 231,
					},
				}
				walker := &walkertype.Walker{}

				Expect(walker.WalkMap(input, func(options *walkertype.WalkerOnWalkMapOptions) *walkertype.WalkerReturn {
					if options.Kind == reflect.String {
						options.Document[options.Key] = options.Value.(string) + "added"

						return &walkertype.WalkerReturn{
							Handled: true,
						}
					}

					return &walkertype.WalkerReturn{
						Handled: false,
					}
				})).To(Equal(walkertype.TypeMap{
					"key1": walkertype.TypeMap{
						"key2": walkertype.TypeMap{
							"key3": walkertype.TypeArray{"arrayval1", "arrayval2", "arrayval3"},
						},
						"key4": "string1added",
						"key5": 231,
					},
				}))
			})
		})
	})
})
