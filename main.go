// Package configuration parsed toml files in specific order to generate a validated struct
package configuration

import (
	os "os"
	path "path"

	common "github.com/beckend/go-config/pkg/common"
	singletons "github.com/beckend/go-config/pkg/singletons"
	validation "github.com/beckend/go-config/pkg/validation"

	spew "github.com/davecgh/go-spew/spew"
	color "github.com/fatih/color"
	config "github.com/gookit/config/v2"
	toml "github.com/gookit/config/v2/toml"
	conditional "github.com/mileusna/conditional"
)

// CallbackGetConfigOptions callback options
type CallbackGetConfigOptions struct {
	Singletons  singletons.Singletons
	Config      config.Config
	FailOnError common.FailOnErrorFunc
	LogSpew     CallbackGeneric
	Validate    CallbackValidate
}

// GetConfigOptions GetConfig options
type GetConfigOptions struct {
	CreateConfig CallbackGetConfig
	EnvKeyRunEnv string
	PathConfigs  string
}

// GetEnv gets environment with fallback
func GetEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

// GetConfig read configurations with priority, checking if the files exists or not
func GetConfig(options GetConfigOptions) interface{} {
	instanceConfig := config.NewEmpty("main-configuration")
	instanceConfig.WithOptions(config.ParseEnv)
	instanceConfig.AddDriver(toml.Driver)
	_, envKeyUserExists := os.LookupEnv(options.EnvKeyRunEnv)
	envRun := GetEnv(conditional.String(envKeyUserExists, options.EnvKeyRunEnv, "RUN_ENV"), "development")

	for _, pathFile := range [...](string){
		// the order to load is base, env specific, then local, where the next overrides the previous values
		path.Join(options.PathConfigs, "base.toml"),
		path.Join(options.PathConfigs, envRun+".toml"),
		path.Join(options.PathConfigs, "local.toml"),
	} {
		if _, err := os.Stat(pathFile); err == nil {
			common.FailOnError(instanceConfig.LoadFiles(pathFile))
		}
	}

	singletonsInstance := singletons.GetSingletons()

	return options.CreateConfig(CallbackGetConfigOptions{
		Singletons: *singletonsInstance,

		Config: *instanceConfig,

		FailOnError: common.FailOnError,

		LogSpew: func(x ...interface{}) {
			color.New(color.FgHiYellow).Println("Parsed environment variables:")
			color.Set(color.FgHiBlue)
			spew.Dump(x...)
			color.Unset()
		},

		Validate: func(x interface{}) {
			singletonsInstance.Validation.Utils.ValidateStruct(validation.ValidatorUtilsValidateStructOptions{
				PrefixError:  "Environment variable error - ",
				Validate:     singletonsInstance.Validation.Validate,
				TheStruct:    x,
				PanicOnError: true,
			})
		},
	})
}

type (
	// CallbackGetConfig type to be used in struct
	CallbackGetConfig func(options CallbackGetConfigOptions) interface{}
	// CallbackGeneric type to be used in struct
	CallbackGeneric func(x ...interface{})
	// CallbackValidate type to be used in struct
	CallbackValidate func(interface{})
)
