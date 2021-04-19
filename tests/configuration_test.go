package configuration_test

import (
	os "os"
	path "path"
	runtime "runtime"
	testing "testing"

	configuration "github.com/beckend/go-config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestValidateStructOne struct {
	AccessKey string `validate:"required"`
	RunEnV    string `validate:"required"`
}

type TestValidateStructOneFail struct {
	AccessKey2 string `validate:"required"`
	RunEnV     string `validate:"required"`
}

func TestPkgAnagram(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "configuration Suite")
}

var _ = Describe("configuration", func() {
	_, pathCurrentFile, _, _ := runtime.Caller(0)
	pathFixtures := path.Join(pathCurrentFile, "../fixtures")

	Context("GetConfig", func() {
		When("only base.toml is present", func() {
			It("values are set", func() {
				result := configuration.GetConfig(configuration.GetConfigOptions{
					CreateConfig: func(options configuration.CallbackGetConfigOptions) interface{} {
						returned := TestValidateStructOne{}

						options.FailOnError(options.Config.BindStruct("", &returned))
						options.Validate(returned)

						return returned
					},
					EnvKeyRunEnv: "RUN_ENV",
					PathConfigs:  path.Join(pathFixtures, "configs-base"),
				}).(TestValidateStructOne)

				Expect(result.RunEnV).To(Equal("development"))
				Expect(result.AccessKey).To(Equal("AccessKey"))
			})

			It("panics upon validation failure", func() {
				Expect(func() {
					configuration.GetConfig(configuration.GetConfigOptions{
						CreateConfig: func(options configuration.CallbackGetConfigOptions) interface{} {
							returned := TestValidateStructOneFail{}

							options.FailOnError(options.Config.BindStruct("", &returned))
							options.Validate(returned)

							return returned
						},
						EnvKeyRunEnv: "RUN_ENV",
						PathConfigs:  path.Join(pathFixtures, "configs-base"),
					})
				}).To(PanicWith(MatchRegexp(`.+validation failed`)))
			})
		})

		When("base.toml is present, env.toml is present", func() {
			It("uses env key EnvKeyRunEnv provided because [env].toml exist and overwrites base.toml", func() {
				keyEnvTarget := "RUN_ENV_CUSTOM"
				keyEnvTargetValue := "staging"
				os.Setenv(keyEnvTarget, keyEnvTargetValue)
				defer os.Unsetenv(keyEnvTarget)

				result := configuration.GetConfig(configuration.GetConfigOptions{
					CreateConfig: func(options configuration.CallbackGetConfigOptions) interface{} {
						returned := TestValidateStructOne{}

						options.FailOnError(options.Config.BindStruct("", &returned))
						options.Validate(returned)

						return returned
					},
					EnvKeyRunEnv: "RUN_ENV_CUSTOM",
					PathConfigs:  path.Join(pathFixtures, "configs-env"),
				}).(TestValidateStructOne)

				Expect(result.RunEnV).To(Equal(keyEnvTargetValue))
			})
		})

		When("base.toml is present, env.toml is present, local.toml is present, and using default EnKeyRunEnv", func() {
			It("local.toml has the last word.", func() {
				result := configuration.GetConfig(configuration.GetConfigOptions{
					CreateConfig: func(options configuration.CallbackGetConfigOptions) interface{} {
						returned := TestValidateStructOne{}

						options.FailOnError(options.Config.BindStruct("", &returned))
						options.Validate(returned)

						return returned
					},
					PathConfigs: path.Join(pathFixtures, "configs-local"),
				}).(TestValidateStructOne)

				Expect(result.RunEnV).To(Equal("local-overwritten"))
			})
		})
	})
})
