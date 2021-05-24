package config_test

import (
	os "os"
	path "path"
	filepath "path/filepath"
	runtime "runtime"
	strings "strings"

	config "github.com/beckend/go-config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestMapstructure struct {
	Key1 string `mapstructure:"AccessKey" validate:"required"`
}

type TestValidateStructOne struct {
	AccessKey string `validate:"required"`
	RunEnV    string `validate:"required"`
	Shell     string `validate:"required"`
	Password  string `validate:"required"`
}

type TestValidateStructTwo struct {
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

	Context("New", func() {
		When("mapstructure annotations are present", func() {
			It("allows for modifications before validating the struct", func() {
				var result TestMapstructure
				_, err := config.New(&config.NewOptions{
					ConfigUnmarshal: &result,
					PathConfigs:     path.Join(pathFixtures, "configs-base"),
				})
				if err != nil {
					panic(err)
				}

				Expect(result.Key1).To(Equal("AccessKey"))
			})
		})

		When("option LoadConfigs", func() {
			It("allows user overrides and keeps what is not overriden", func() {
				var result TestValidateStructOne
				_, err := config.New(&config.NewOptions{
					ConfigUnmarshal: &result,
					EnvKeyRunEnv:    "RUN_ENV",
					LoadConfigs: func(options *config.LoadConfigsOptions) ([][]byte, error) {
						b1, err := options.TOML.BytesToJSON([]byte("RunEnv = 'overriden'"))
						if err != nil {
							return nil, err
						}

						b2, err := options.TOML.StringToJSON("AccessKey = 'overriden'")
						if err != nil {
							return nil, err
						}

						b3, err := options.TOML.ReaderToJSON(strings.NewReader("Shell = 'overriden'"))
						if err != nil {
							return nil, err
						}

						return [][]byte{b1, b2, b3}, nil
					},
					PathConfigs: path.Join(pathFixtures, "configs-base"),
				})
				if err != nil {
					panic(err)
				}

				Expect(result.RunEnV).To(Equal("overriden"))
				Expect(result.AccessKey).To(Equal("overriden"))
				Expect(result.Shell).To(Equal("overriden"))
				Expect(result.Password).To(Equal("defaultpassword"))
			})
		})

		When("option OnConfigBeforeValidation", func() {
			It("allows for modifications before validating the struct", func() {
				var result TestValidateStructOne
				_, err := config.New(&config.NewOptions{
					ConfigUnmarshal: &result,
					EnvKeyRunEnv:    "RUN_ENV",
					OnConfigBeforeValidation: func(options *config.OnConfigBeforeValidationOptions) error {
						myConfig := options.ConfigUnmarshal.(*TestValidateStructOne)
						myConfig.Password = "nope"
						return nil
					},
					PathConfigs: path.Join(pathFixtures, "configs-base"),
				})
				if err != nil {
					panic(err)
				}

				Expect(result.Password).To(Equal("nope"))
			})
		})

		When("only base.toml is present", func() {
			It("values are set", func() {
				var result TestValidateStructOne
				_, err := config.New(&config.NewOptions{
					ConfigUnmarshal: &result,
					EnvKeyRunEnv:    "RUN_ENV",
					PathConfigs:     path.Join(pathFixtures, "configs-base"),
				})
				if err != nil {
					panic(err)
				}

				Expect(result.RunEnV).To(Equal("development"))
				Expect(result.AccessKey).To(Equal("AccessKey"))
				Expect(result.Password).To(Equal("defaultpassword"))
				Expect(result.Shell).ToNot(Equal("${SHELL}"))
			})

			It("panics upon validation failure", func() {
				Expect(func() {
					var input TestValidateStructOneFail

					_, err := config.New(&config.NewOptions{
						ConfigUnmarshal: &input,
						EnvKeyRunEnv:    "RUN_ENV",
						PathConfigs:     path.Join(pathFixtures, "configs-base"),
					})

					Expect(err.Error()).To(ContainSubstring("validation failed"))

					panic(err)
				}).To(Panic())
			})

			When("base.toml is present, env.toml is present", func() {
				It("uses env key EnvKeyRunEnv provided because [env].toml exist and overwrites base.toml", func() {
					var result TestValidateStructTwo
					keyEnvTarget := "RUN_ENV_CUSTOM"
					keyEnvTargetValue := "staging"
					os.Setenv(keyEnvTarget, keyEnvTargetValue)
					defer os.Unsetenv(keyEnvTarget)

					_, err := config.New(&config.NewOptions{
						ConfigUnmarshal: &result,
						EnvKeyRunEnv:    keyEnvTarget,
						PathConfigs:     path.Join(pathFixtures, "configs-env"),
					})
					if err != nil {
						panic(err)
					}

					Expect(result.RunEnV).To(Equal(keyEnvTargetValue))
				})
			})

			When("base.toml is present, env.toml is present, local.toml is present, and using default EnKeyRunEnv", func() {
				It("local.toml has the last word.", func() {
					var result TestValidateStructOne
					_, err := config.New(&config.NewOptions{
						ConfigUnmarshal: &result,
						PathConfigs:     path.Join(pathFixtures, "configs-local"),
					})
					if err != nil {
						panic(err)
					}

					Expect(result.RunEnV).To(Equal("local-overwritten"))
				})
			})
		})
	})
})
