package walkertype_test

import (
	fmt "fmt"
	path "path"
	filepath "path/filepath"
	reflect "reflect"
	runtime "runtime"
	testing "testing"

	common "github.com/beckend/go-config/pkg/common"
	file "github.com/beckend/go-config/pkg/file"
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
	_, pathCurrentFile, _, _ := runtime.Caller(0)
	pathDirCurrent, _ := filepath.Split(pathCurrentFile)
	pathFixtures := path.Join(pathDirCurrent, "fixtures")

	Context("Walk", func() {
		When("advanced tests", func() {
			It("does not panic", func() {
				Expect(func() {
					input, err := file.JSONFileToMap(path.Join(pathFixtures, "test1.json"))
					common.FailOnError(err)

					res := walkertype.Walk(&walkertype.WalkOptions{
						Object: input,
						OnKind: func(oosvo *walkertype.OnKindOptions) *walkertype.OnKindWalkReturn {
							if oosvo.CaseKind == reflect.String {
								str := oosvo.Original.String()
								if str == "something" {
									oosvo.Copy.SetString(str + " more")
								}

								return &walkertype.OnKindWalkReturn{
									Handled: true,
								}
							}

							return &walkertype.OnKindWalkReturn{
								Handled: false,
							}
						},

						OnMapKey: func(omko *walkertype.OnMapKeyOptions) *walkertype.OnKindWalkReturn {
							if omko.Key.Kind() == reflect.String {
								valueStr := omko.Key.String()
								if valueStr == "RunEnv" {
									// does not retain the property value
									omko.Copy.SetMapIndex(reflect.ValueOf(fmt.Sprintf("%s is production", valueStr)), omko.ValueCopy)

									return &walkertype.OnKindWalkReturn{
										Handled: true,
									}
								}
							}

							return &walkertype.OnKindWalkReturn{
								Handled: false,
							}
						},
					})

					result := res.(map[string]interface{})
					Expect(result["NameService"]).To(Equal("something more"))
					// @TODO fix this, need to figure out how to modify the map with it's property value intact
					// Expect(result["RunEnv is production"]).To(Equal("development"))
				}).ToNot(Panic())
			})
		})
	})
})
