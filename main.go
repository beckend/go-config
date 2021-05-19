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

// CallbackNewOptions callback options
type CallbackNewOptions struct {
	Singletons  singletons.Singletons
	Config      config.Config
	FailOnError common.FailOnErrorFunc
	LogSpew     CallbackGeneric
	Validate    CallbackValidate
}

// NewOptions GetConfig options
type NewOptions struct {
	CreateConfig CallNewConfig
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

// New read configurations with priority, checking if the files exists or not
func New(options NewOptions) interface{} {
	instanceConfig := config.NewEmpty("main-configuration")
	instanceConfig.AddDriver(toml.Driver)
	instanceConfig.WithOptions(config.ParseEnv)

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

	singletonsInstance := singletons.New()

	return options.CreateConfig(CallbackNewOptions{
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
			singletonsInstance.Validation.ValidateStruct(validation.ValidatorUtilsValidateStructOptions{
				PrefixError:  "Environment variable error - ",
				TheStruct:    x,
				PanicOnError: true,
			})
		},
	})
}

type (
	// CallNewConfig type to be used in struct
	CallNewConfig func(options CallbackNewOptions) interface{}
	// CallbackGeneric type to be used in struct
	CallbackGeneric func(x ...interface{})
	// CallbackValidate type to be used in struct
	CallbackValidate func(interface{})
)
