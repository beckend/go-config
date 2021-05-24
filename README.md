# go-config parses config files in specific order to generate a validated struct

## Install

```shell
go get -u github.com/beckend/go-config
```

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
