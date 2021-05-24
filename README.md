# go-config parses config files in specific order to generate a validated struct

## Install

```shell
go get -u github.com/beckend/go-config
```

---

### Description and usage

- a directory with config files
- the directory of files with be parsed in order: `base.toml` -> `[env].toml` -> `local.toml` and the later one will overwrite the values of the previous, all files are optionally existing.
- `[env].toml` is calulated by providing in options struct with key `EnvKeyRunEnv`, so set for example to `RUN_ENV` it will read environment variable `RUN_ENV` and if the value is for example `staging` it becomes `staging.toml`.
- the config structs which are to be validated can be annotated correctly using struct tags according to https://github.com/go-playground/validator - which is being used in this library.
- environment variable substitution, any value which looks like `"${MY_VAR}"` will be replaced by environment variables,
  default is also supported when env variable is missing => `"${MY_VAR|defaultValue}"`
- stuct tags from https://github.com/mitchellh/mapstructure works, uses decode internally after json.Unmarshal(), see https://github.com/mitchellh/mapstructure/blob/master/mapstructure_examples_test.go

Example:

when `/path/to/directory/with-configs` contains `base.toml` with contents:

```toml
APIKeyGithub = 'secret-key'
```

`local.toml` with contents:

```toml
APIKeyGithub = 'secret-key-local-dev'
passwordfromenv = '${PASSWORD}'
username = '${USERNAME|nobody}'
```

```golang
package mypackage

import (
  fmt "fmt"
  config "github.com/beckend/go-config"
  path "path"
)

// See https://github.com/go-playground/validator
type MyConfig struct {
  APIKeyGithub string `validate:"required"`
  PasswordFromEnv string `mapstructure:"passwordfromenv" validate:"required"`
  UserName string `mapstructure:"username" validate:"required"`
}

func main() {
  var result MyConfig
  _, err := config.New(&config.NewOptions{
    ConfigUnmarshal: &result,
    EnvKeyRunEnv:    "RUN_ENV",
    PathConfigs:     path.Join("/my/directory-with-configs", "configs-base"),
  })
  if err != nil {
    panic(err)
  }

  // prints secret-key-local-dev since local.toml is the last parsed in priority chain
  fmt.Println(myConfig.APIKeyGithub)
  // Whatever environment PASSWORD was set to
  fmt.Println(myConfig.PasswordFromEnv)
  // nobody if USERNAME is unset, otherwise the value of existing environment variable
  fmt.Println(myConfig.UserName)
}
```

---

### option `LoadConfigs` allows loading from custom sources

Returning an array of json marshalled bytes, the order matter where the later one will override the previous.
See `main_test.go` for details, the gist is

`base.toml`

```toml
RunEnv = 'development'
AccessKey = "AccessKey"
Password = "${____password___|defaultpassword}"
Shell = "${SHELL}"
```

```golang
type TestValidateStructOne struct {
	AccessKey string `validate:"required"`
	RunEnV    string `validate:"required"`
	Shell     string `validate:"required"`
	Password  string `validate:"required"`
}

err := os.Setenv("RUN_ENV", "staging")
if err !=nil {
  panic(err)
}

var result TestValidateStructOne
_, err := config.New(&config.NewOptions{
  ConfigUnmarshal: &result,
  EnvKeyRunEnv:    "RUN_ENV",
  LoadConfigs: func(options *config.LoadConfigsOptions) ([][]byte, error) {
    Expect(options.RunEnv).To(Equal("staging"))

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
```

---

### option `OnConfigBeforeValidation` allows modifications before struct is going to be validated to do custom logic before validation

Good place to add complex logic to read/replace variables, at this stage all env variables have been replaced

`base.toml`

```toml
RunEnv = 'development'
AccessKey = "AccessKey"
Password = "${____password___|defaultpassword}"
Shell = "${SHELL}"
```

```golang
type TestValidateStructOne struct {
	AccessKey string `validate:"required"`
	RunEnV    string `validate:"required"`
	Shell     string `validate:"required"`
	Password  string `validate:"required"`
}

var result TestValidateStructOne
_, err := config.New(&config.NewOptions{
  ConfigUnmarshal: &result,
  EnvKeyRunEnv:    "RUN_ENV",
  OnConfigBeforeValidation: func(options *config.OnConfigBeforeValidationOptions) error {
    myConfig := options.ConfigUnmarshal.(*TestValidateStructOne)
    myConfig.Password = "nope"
    // /bin/*** depends on your environment
    fmt.Println(myConfig.Shell)
    return nil
  },
  PathConfigs: path.Join(pathFixtures, "configs-base"),
})
if err != nil {
  panic(err)
}

Expect(result.Password).To(Equal("nope"))
```
