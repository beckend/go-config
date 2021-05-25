package config_test

import (
	errors "errors"
	fs "io/fs"
	os "os"
	path "path"
	filepath "path/filepath"
	runtime "runtime"
	strings "strings"

	config "github.com/beckend/go-config"
	common "github.com/beckend/go-config/pkg/common"

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

var _ = Describe("pkg main", func() {
	_, pathCurrentFile, _, _ := runtime.Caller(0)
	pathDirCurrent, _ := filepath.Split(pathCurrentFile)
	pathFixtures := path.Join(pathDirCurrent, "tests/fixtures")

	Context("New", func() {
		When("mapstructure annotations are present", func() {
			It("renames field test", func() {
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

		When("RunEnv env is missing", func() {
			It("string will be empty", func() {
				var result TestValidateStructTwo
				_, err := config.New(&config.NewOptions{
					ConfigUnmarshal: &result,
					PathConfigs:     path.Join(pathFixtures, "does/not-exist"),
					OnConfigBeforeValidation: func(options *config.OnConfigBeforeValidationOptions) error {
						cfg := options.ConfigUnmarshal.(*TestValidateStructTwo)
						cfg.AccessKey = "somethings"

						if cfg.RunEnV != "" {
							return errors.New("this is not supposed to happen")
						}

						return nil
					},
				})

				Expect(err.Error()).To(ContainSubstring("struct validation failed"))
			})
		})

		When("option LoadConfigs", func() {
			It("allows user overrides and keeps what is not overriden", func() {
				keyEnvTarget := "RUN_ENV_CUSTOM"
				keyEnvTargetValue := "staging"
				err := os.Setenv(keyEnvTarget, keyEnvTargetValue)
				defer os.Unsetenv(keyEnvTarget)
				if err != nil {
					panic(err)
				}

				var result TestValidateStructOne
				_, err = config.New(&config.NewOptions{
					ConfigUnmarshal: &result,
					EnvKeyRunEnv:    keyEnvTarget,
					LoadConfigs: func(options *config.LoadConfigsOptions) ([][]byte, error) {
						Expect(options.RunEnv).To(Equal(keyEnvTargetValue))

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
					err := os.Setenv(keyEnvTarget, keyEnvTargetValue)
					common.FailOnError(err)
					defer os.Unsetenv(keyEnvTarget)

					_, err = config.New(&config.NewOptions{
						ConfigUnmarshal: &result,
						EnvKeyRunEnv:    keyEnvTarget,
						PathConfigs:     path.Join(pathFixtures, "configs-env"),
					})
					common.FailOnError(err)

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
					common.FailOnError(err)

					Expect(result.RunEnV).To(Equal("local-overwritten"))
				})
			})

			When("configs-custom test1.toml", func() {
				It("works", func() {
					var result map[string]interface{}
					_, err := config.New(&config.NewOptions{
						ConfigUnmarshal: &result,
						LoadConfigs: func(options *config.LoadConfigsOptions) ([][]byte, error) {
							b1, err := options.TOML.FileReaderCallbackToJSON(func() (fs.File, error) {
								return os.Open(path.Join(pathFixtures, "configs-custom/test1.toml"))
							})
							common.FailOnError(err)

							return [][]byte{b1}, nil
						},
					})
					common.FailOnError(err)

					Expect(result["RunEnv"]).To(Equal("development"))
					Expect(result["Shell"]).ToNot(Equal("${SHELL}"))
				})
			})
		})
	})
})
