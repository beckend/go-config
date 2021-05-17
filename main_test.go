package configuration_test

import (
	os "os"
	path "path"
	filepath "path/filepath"
	runtime "runtime"

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

var _ = Describe("pkg validation", func() {
	_, pathCurrentFile, _, _ := runtime.Caller(0)
	pathDirCurrent, _ := filepath.Split(pathCurrentFile)
	pathFixtures := path.Join(pathDirCurrent, "tests/fixtures")

	Context("GetEnv", func() {
		When("env exists", func() {
			It("returns correct value", func() {
				Expect(len(configuration.GetEnv("SHELL", "")) > 0).To(Equal(true))
			})
		})

		When("env does not exist", func() {
			It("returns fallback", func() {
				Expect(configuration.GetEnv("FDSCCVB##csdas#@!CS", "")).To(Equal(""))
			})
		})
	})

	Context("New", func() {
		When("only base.toml is present", func() {
			It("values are set", func() {
				result := configuration.New(configuration.NewOptions{
					CreateConfig: func(options configuration.CallbackNewOptions) interface{} {
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
					configuration.New(configuration.NewOptions{
						CreateConfig: func(options configuration.CallbackNewOptions) interface{} {
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

				result := configuration.New(configuration.NewOptions{
					CreateConfig: func(options configuration.CallbackNewOptions) interface{} {
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
				result := configuration.New(configuration.NewOptions{
					CreateConfig: func(options configuration.CallbackNewOptions) interface{} {
						returned := TestValidateStructOne{}

						options.FailOnError(options.Config.BindStruct("", &returned))
						options.Validate(returned)

						options.LogSpew(returned)

						return returned
					},
					PathConfigs: path.Join(pathFixtures, "configs-local"),
				}).(TestValidateStructOne)

				Expect(result.RunEnV).To(Equal("local-overwritten"))
			})
		})
	})
})
