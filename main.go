// Package config parsed toml files in specific order to generate a validated struct
package config

import (
	errors "errors"
	io "io"
	os "os"
	path "path"
	reflect "reflect"

	environment "github.com/beckend/go-config/pkg/environment"
	file "github.com/beckend/go-config/pkg/file"
	singletons "github.com/beckend/go-config/pkg/singletons"
	validation "github.com/beckend/go-config/pkg/validation"
	walkertype "github.com/beckend/go-config/pkg/walker-type"
	validator "github.com/go-playground/validator/v10"
	envutil "github.com/gookit/goutil/envutil"
	jsoniter "github.com/json-iterator/go"

	conditional "github.com/mileusna/conditional"
	mapstructure "github.com/mitchellh/mapstructure"
)

var (
	defaultRUNENV = "development"
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Config struct {
	ErrorsValidation *validator.ValidationErrors
}

type OnConfigBeforeValidationOptions struct {
	ConfigUnmarshal interface{}
}

type LoadConfigsOptionsTOML struct {
	FileToJSON   func(string) ([]byte, error)
	StringToJSON func(string) ([]byte, error)
	BytesToJSON  func([]byte) ([]byte, error)
	ReaderToJSON func(io.Reader) ([]byte, error)
}

type LoadConfigsOptions struct {
	TOML   *LoadConfigsOptionsTOML
	RunEnv string
}

type (
	// Allows user to have a shot on the config before it validats
	OnConfigBeforeValidation func(options *OnConfigBeforeValidationOptions) error
	// Allow user to load configs, user has to return a []byte which has been through json.Marshal into []byte
	// The payload is in the end going to be json.Unmarshaled
	LoadConfigs func(options *LoadConfigsOptions) ([][]byte, error)
)

type NewOptions struct {
	OnConfigBeforeValidation OnConfigBeforeValidation
	LoadConfigs              LoadConfigs
	EnvKeyRunEnv             string
	PathConfigs              string
	ConfigUnmarshal          interface{}
}

// New read configurations with priority, the later overrides the previous
func New(options *NewOptions) (*Config, error) {
	_, envKeyUserExists := os.LookupEnv(options.EnvKeyRunEnv)
	envRun := environment.GetEnv(conditional.String(envKeyUserExists, options.EnvKeyRunEnv, "RUN_ENV"), defaultRUNENV)
	var filesToBeMerged []string

	for _, pathFile := range [...](string){
		// the order to load is base, env specific, then local, where the next overrides the previous values
		path.Join(options.PathConfigs, "base.toml"),
		path.Join(options.PathConfigs, envRun+".toml"),
		path.Join(options.PathConfigs, "local.toml"),
	} {
		if _, err := os.Stat(pathFile); err == nil {
			filesToBeMerged = append(filesToBeMerged, pathFile)
		}
	}

	bytesJSONMerged, err := file.TOMLFilesToMergedJSON(filesToBeMerged)
	if err != nil {
		return nil, err
	}

	if options.LoadConfigs != nil {
		byteSlicesUser, err := options.LoadConfigs(&LoadConfigsOptions{
			TOML: &LoadConfigsOptionsTOML{
				FileToJSON:   file.TOMLFileToJSON,
				BytesToJSON:  file.TOMLBytesToJSON,
				StringToJSON: file.TOMLStringToJSON,
				ReaderToJSON: file.TOMLReaderToJSON,
			},
			RunEnv: envRun,
		})
		if err != nil {
			return nil, err
		}

		// prepend bytesJSONMerged into byteSlicesUser so the user files overrides the initial config
		byteSlicesUser = append([][]byte{bytesJSONMerged}, byteSlicesUser...)
		bytesJSONMerged, err = file.TOMLBytesToMergedJSON(byteSlicesUser)
		if err != nil {
			return nil, err
		}
	}

	// convert to a generic map interface to replace env variables
	var configMap walkertype.TypeMap
	err = json.Unmarshal(bytesJSONMerged, &configMap)
	if err != nil {
		return nil, err
	}

	(&walkertype.Walker{}).WalkMap(configMap, func(options *walkertype.WalkerOnWalkMapOptions) *walkertype.WalkerReturn {
		if options.Kind == reflect.String {
			options.Document[options.Key] = envutil.ParseEnvValue(options.Value.(string))

			return &walkertype.WalkerReturn{
				Handled: true,
			}
		}

		return &walkertype.WalkerReturn{
			Handled: false,
		}
	})

	mapstructure.Decode(configMap, &options.ConfigUnmarshal)

	if options.OnConfigBeforeValidation != nil {
		err = options.OnConfigBeforeValidation(&OnConfigBeforeValidationOptions{
			ConfigUnmarshal: options.ConfigUnmarshal,
		})

		if err != nil {
			return nil, err
		}
	}

	errsValidation := singletons.New().Validation.ValidateStruct(validation.ValidatorUtilsValidateStructOptions{
		PrefixError:  "Config struct validation error - ",
		TheStruct:    options.ConfigUnmarshal,
		PanicOnError: false,
	})

	if errsValidation != nil && len(*errsValidation) > 0 {
		err = errors.New("config struct validation failed")
	} else {
		err = nil
	}

	return &Config{
		ErrorsValidation: errsValidation,
	}, err
}
